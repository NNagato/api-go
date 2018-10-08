package core

import (
	"errors"
	"log"
	"strings"

	"github.com/KyberNetwork/api-server/common"
	"github.com/KyberNetwork/api-server/fetcher"
	"github.com/KyberNetwork/api-server/storage"
)

const (
	ETH_TO_WEI = 1000000000000000000
)

type Core struct {
	fetcher *fetcher.Fetcher
	storage *storage.RamStorage
}

func NewCore(fetcher *fetcher.Fetcher, storage *storage.RamStorage) *Core {
	return &Core{
		fetcher: fetcher,
		storage: storage,
	}
}

func (self *Core) GetListTokenAPI() ([]common.Token, error) {
	return self.fetcher.GetListTokenAPI()
}

func (self *Core) GetRateBuy(idArr []string, qtyArr [][]float64) ([]common.RateBuy, error) {
	infoData := self.fetcher.GetInfoData()
	listTokens := self.fetcher.GetMapListToken()
	rates := self.storage.GetRate()
	rateBuy := []common.RateBuy{}
	ethAddr := infoData.ETHAddr
	ethSymbol := infoData.ETHSymbol

	dataGetRate, err := GetDataGetRate(ethAddr, ethSymbol, idArr, qtyArr, listTokens, rates)
	if err != nil {
		log.Println(err)
		return rateBuy, err
	}

	trueRate, err := self.fetcher.GetRate([]common.Rate{}, &dataGetRate)
	if err != nil {
		log.Println(err)
		return rateBuy, err
	}

	if len(dataGetRate.SourceArr) != len(*trueRate) {
		return rateBuy, errors.New("Can't get rate")
	}

	for index, id := range idArr {
		idLower := strings.ToLower(id)
		if token, ok := listTokens[idLower]; ok {
			dstSymbol := token.Symbol
			dstQty := qtyArr[index]
			srcQty, err := GetBuyRate(*trueRate, dstQty, dstSymbol)
			if err != nil {
				log.Println(err)
				return rateBuy, err
			}
			rateBuy = append(rateBuy, common.RateBuy{
				SrcID:  ethAddr,
				DstID:  id,
				SrcQty: srcQty,
				DstQty: dstQty,
			})
		}
	}
	return rateBuy, nil
}

func (self *Core) GetIsNewRate() bool {
	return self.storage.GetIsNewRate()
}
