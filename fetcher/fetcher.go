package fetcher

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"crypto/hmac"
	"crypto/sha512"

	"github.com/KyberNetwork/api-server/cmd/config"
	"github.com/KyberNetwork/api-server/common"
	fetchercore "github.com/KyberNetwork/api-server/fetcher/fetcher-core"
	ethereum "github.com/ethereum/go-ethereum/common"
)

type Fetcher struct {
	mu             *sync.RWMutex
	fetcherIns     []FetcherInterfce
	infoData       *config.Config
	contractHanler *EthereumHandler
	tokens         []common.Token

	errUpdateToken error
}

func (self *Fetcher) GetInfoData() *config.Config {
	infoData := self.infoData
	listToken := self.GetListToken()
	infoData.Tokens = listToken
	return infoData
}

// get list token array
func (self *Fetcher) GetListToken() []common.Token {
	self.mu.Lock()
	defer self.mu.Unlock()
	return self.tokens
}

// get list token (a map with key is token's address)
func (self *Fetcher) GetMapListToken() map[string]common.Token {
	arrayToken := self.GetListToken()
	mapListToken := make(map[string]common.Token)
	for _, token := range arrayToken {
		mapListToken[strings.ToLower(token.Address)] = token
	}
	return mapListToken
}

// get token include error
func (self *Fetcher) GetListTokenForAPI() ([]common.Token, error) {
	self.mu.Lock()
	defer self.mu.Unlock()
	return self.tokens, self.errUpdateToken
}

// Update list tokens from core
func (self *Fetcher) UpdateListToken(listToken []common.Token, err error) {
	self.mu.Lock()
	defer self.mu.Unlock()
	if err == nil {
		self.tokens = listToken
	}
	self.errUpdateToken = err
}

func NewFetcher(infoData *config.Config) (*Fetcher, error) {
	ethereum, err := NewEthereumHandler(infoData.Network, infoData.Wrapper)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	fetcherIns := []FetcherInterfce{}

	for _, endpoint := range infoData.NodeEndpoint {
		bcFetcher, err := fetchercore.NewBlockchainFetcher(endpoint, infoData.Network, infoData.Wrapper)
		if err != nil {
			log.Println(err)
			continue
		}
		fetcherIns = append(fetcherIns, bcFetcher)
	}
	etherscanFetcher := fetchercore.NewEtherScanFetcher(infoData.EtherScanAPI)
	fetcherIns = append(fetcherIns, etherscanFetcher)
	mu := &sync.RWMutex{}

	fetcher := &Fetcher{
		mu:             mu,
		infoData:       infoData,
		contractHanler: ethereum,
		fetcherIns:     fetcherIns,
	}

	err = fetcher.RunUpdateListToken()
	if err != nil {
		log.Println(err)
		fetcher.errUpdateToken = err
	}

	// interval fetch listToken
	ticker := time.NewTicker(3600 * time.Second)
	go func() {
		for {
			<-ticker.C
			err = fetcher.RunUpdateListToken()
			if err != nil {
				log.Println(err)
			}
		}
	}()

	return fetcher, nil
}

// Update listoken
func (self *Fetcher) RunUpdateListToken() error {
	infoData := self.GetInfoData()
	secretCore := infoData.SecretCore
	nonce := time.Now().UTC().UnixNano() / int64(1000000)
	message := fmt.Sprintf("nonce=%v", nonce)

	// calculate signed for core API
	h := hmac.New(sha512.New, []byte(secretCore))
	h.Write([]byte(message))
	signedMessage := hex.EncodeToString(h.Sum(nil))

	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/setting/token-settings?%s", self.infoData.CoreEndpoint, message), nil)
	if err != nil {
		log.Println(err)
		self.UpdateListToken([]common.Token{}, err)
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("signed", signedMessage)
	response, err := client.Do(req)
	if err != nil {
		log.Println(err)
		err = errors.New("Cant get data from core")
		self.UpdateListToken([]common.Token{}, err)
		return err
	}

	defer (response.Body).Close()
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Print(err)
		self.UpdateListToken([]common.Token{}, err)
		return err
	}
	listTokenResponse := common.ListTokenResponse{}
	err = json.Unmarshal(b, &listTokenResponse)
	if err != nil {
		log.Print(err)
		self.UpdateListToken([]common.Token{}, err)
		return err
	}
	if listTokenResponse.Success == true {
		self.UpdateListToken(listTokenResponse.Data, nil)
	}
	return nil
}

func (self *Fetcher) GetListTokenAPI() ([]common.Token, error) {
	tokens, err := self.GetListTokenForAPI()
	return tokens, err
}

// ----------------------- fetch blockchain ---------------------------
// read contract

func (self *Fetcher) GetRate(currentRate []common.Rate, customData *common.DataGetRate) (*[]common.Rate, error) {
	infoData := self.GetInfoData()
	var sourceArr, sourceSymbolArr, destArr, destSymbolArr []string
	var amountArr []*big.Int
	if len(currentRate) == 0 && customData != nil {
		sourceArr = customData.SourceArr
		sourceSymbolArr = customData.SourceSymbolArr
		destArr = customData.DstArr
		destSymbolArr = customData.DstSymbolArr
		amountArr = customData.AmountArr
	} else {
		sourceArr, sourceSymbolArr, destArr, destSymbolArr, amountArr = MakeDataGetRate(infoData, currentRate)
	}

	dataAbi, err := self.contractHanler.EncodeRateData(sourceArr, destArr, amountArr)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	for _, fetIns := range self.fetcherIns {
		result, err := fetIns.EthCall(self.infoData.Wrapper, ethereum.ToHex(dataAbi))
		if err != nil {
			log.Print(err)
			continue
		}
		rates, err := self.contractHanler.ExtractRateData(result, sourceSymbolArr, destSymbolArr)
		if err != nil {
			log.Print(err)
			continue
		}
		return rates, nil
	}
	return nil, errors.New("Cannot get rate")
}
