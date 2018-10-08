package core

import (
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/KyberNetwork/api-server/common"
)

func GetBuyRate(rates []common.Rate, qtyArr []float64, symbol string) ([]float64, error) {
	var rateString string
	var err error
	for _, rate := range rates {
		if rate.Dest == symbol {
			rateString = rate.Rate
		}
	}
	rateNumber, done := big.NewInt(0).SetString(rateString, 10)
	if done == false {
		err = fmt.Errorf("Cant get rate %s", symbol)
		return []float64{}, err
	}
	rateFloat := big.NewFloat(0).SetInt(rateNumber)

	oneEthToToken, _ := big.NewFloat(0).Quo(rateFloat, big.NewFloat(ETH_TO_WEI)).Float64()
	buyRate := []float64{}
	for _, qty := range qtyArr {
		buy := qty / oneEthToToken
		buyRate = append(buyRate, buy)
	}
	return buyRate, nil
}

func GetDataGetRate(ethAddr string, ethSymbol string,
	idArr []string,
	qtyArr [][]float64,
	listTokens map[string]common.Token,
	rates []common.Rate) (common.DataGetRate, error) {
	sourceArr := []string{}
	sourceSymbolArr := []string{}
	dstArr := []string{}
	dstSymbolArr := []string{}
	amountArr := []*big.Int{}

	for index, id := range idArr {
		idLower := strings.ToLower(id)
		if token, ok := listTokens[idLower]; ok {
			dstSymbol := token.Symbol
			dstQty := qtyArr[index]
			srcQty, err := GetBuyRate(rates, dstQty, dstSymbol)
			if err != nil {
				log.Println(err)
				return common.DataGetRate{}, err
			}
			for _, ele := range srcQty {
				sourceArr = append(sourceArr, ethAddr)
				sourceSymbolArr = append(sourceSymbolArr, ethSymbol)
				dstArr = append(dstArr, token.Address)
				dstSymbolArr = append(dstSymbolArr, dstSymbol)
				amount := common.GetAmountInWei(ele)
				amountArr = append(amountArr, amount)
			}
		}
	}
	dataGetRate := common.DataGetRate{
		SourceArr:       sourceArr,
		SourceSymbolArr: sourceSymbolArr,
		DstArr:          dstArr,
		DstSymbolArr:    dstSymbolArr,
		AmountArr:       amountArr,
	}
	return dataGetRate, nil
}
