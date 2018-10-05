package fetchercore

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/KyberNetwork/api-server/common"
)

type EthereScanFetcher struct {
	etherscanAPI string
	apiKey       string
}

func NewEtherScanFetcher(etherscanAPI string) *EthereScanFetcher {
	return &EthereScanFetcher{
		etherscanAPI: etherscanAPI,
	}
}

// call contract
func (self *EthereScanFetcher) EthCall(toAddr string, dataABI string) (string, error) {
	url := self.etherscanAPI + "/api?module=proxy&action=eth_call&to=" + toAddr + "&data=" + dataABI + "&tag=latest&apikey=" + self.apiKey
	response, err := http.Get(url)

	if err != nil {
		log.Print(err)
		return "", err
	}
	if response.StatusCode != 200 {
		return "", errors.New("Status code is 200")
	}

	defer (response.Body).Close()
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Print(err)
		return "", err
	}
	result := common.ResultRpc{}
	err = json.Unmarshal(b, &result)
	if err != nil {
		log.Print(err)
		return "", err
	}

	return result.Result, nil
}

// ----------- Get balance in WEI

func (self *EthereScanFetcher) GetBalanceAccount() {
	toAddr := "0xddbd2b932c763ba5b1b7ae3b362eac3e8d40121a" // example account
	url := self.etherscanAPI + "/api?module=account&action=balance&address=" + toAddr + "&tag=latest&apikey=" + self.apiKey
	response, err := http.Get(url)

	if err != nil {
		log.Print(err)
		return
	}
	if response.StatusCode != 200 {
		return
	}

	defer (response.Body).Close()
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Print(err)
		return
	}
	result := common.ResultRpc{}
	err = json.Unmarshal(b, &result)
	if err != nil {
		log.Print(err)
		return
	}
	log.Println("balance: ", result.Result)
}

// ----------- EstimateGas example

func (self *EthereScanFetcher) EstimateGas() {
	//
}
