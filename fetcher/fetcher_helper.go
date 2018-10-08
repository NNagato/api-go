package fetcher

import (
	"math/big"

	"github.com/KyberNetwork/api-server/cmd/config"
	"github.com/KyberNetwork/api-server/common"
)

const (
	ETH_TO_WEI = 1000000000000000000
	MIN_ETH    = 0.001
)

func getAmountTokenWithMinETH(rate *big.Int, decimal int) *big.Int {
	rFloat := big.NewFloat(0).SetInt(rate)
	ethFloat := big.NewFloat(ETH_TO_WEI)
	amoutnToken1ETH := rFloat.Quo(rFloat, ethFloat)
	minAmountWithMinETH := amoutnToken1ETH.Mul(amoutnToken1ETH, big.NewFloat(MIN_ETH))
	decimalWei := big.NewInt(0).Exp(big.NewInt(10), big.NewInt(int64(decimal)), nil)
	amountWithDecimal := big.NewFloat(0).Mul(minAmountWithMinETH, big.NewFloat(0).SetInt(decimalWei))
	amountInt, _ := amountWithDecimal.Int(nil)
	return amountInt
}

func MakeDataGetRate(info *config.Config, currentRate []common.Rate) (
	[]string, []string, []string, []string,
	[]*big.Int,
) {
	sourceAddr := make([]string, 0)
	sourceSymbol := make([]string, 0)
	destAddr := make([]string, 0)
	destSymbol := make([]string, 0)
	amount := make([]*big.Int, 0)
	amountETH := make([]*big.Int, 0)
	ethSymbol := info.ETHSymbol
	ethAddr := info.ETHAddr
	minAmountETH := common.GetAmountInWei(MIN_ETH)
	listToken := make(map[string]common.Token)
	for _, token := range info.Tokens {
		listToken[token.Symbol] = token
	}

	if len(currentRate) > 0 {
		for _, rate := range currentRate {
			if rate.Source == "ETH" && rate.Dest == "ETH" {
				continue
			}
			if rate.Source == "ETH" {
				amountToken := big.NewInt(0)
				r := big.NewInt(0)
				r.SetString(rate.Rate, 10)
				destSym := rate.Dest
				decimal := listToken[destSym].Decimals
				if decimal != 0 || r.Cmp(amountToken) != 0 {
					amountToken = getAmountTokenWithMinETH(r, decimal)
				}
				amount = append(amount, amountToken)
				tokenAddr := listToken[destSym].Address
				sourceAddr = append(sourceAddr, tokenAddr)
				sourceSymbol = append(sourceSymbol, destSym)
			} else {
				destAddr = append(destAddr, ethAddr)
				destSymbol = append(destSymbol, ethSymbol)
				amountETH = append(amountETH, minAmountETH)
			}
		}
		sourceAddr = append(sourceAddr, ethAddr)
		destAddr = append(destAddr, ethAddr)
		sourceSymbol = append(sourceSymbol, ethSymbol)
		destSymbol = append(destSymbol, ethSymbol)
		amount = append(amount, minAmountETH)
		amountETH = append(amountETH, minAmountETH)
	} else {
		for _, token := range info.Tokens {
			sourceAddr = append(sourceAddr, token.Address)
			sourceSymbol = append(sourceSymbol, token.Symbol)
			destAddr = append(destAddr, ethAddr)
			destSymbol = append(destSymbol, ethSymbol)
			amount = append(amount, big.NewInt(0))
			amountETH = append(amountETH, minAmountETH)
		}
	}

	sourceArr := append(sourceAddr, destAddr...)
	sourceSymbolArr := append(sourceSymbol, destSymbol...)
	destArr := append(destAddr, sourceAddr...)
	destSymbolArr := append(destSymbol, sourceSymbol...)
	amountArr := append(amount, amountETH...)

	return sourceArr, sourceSymbolArr, destArr, destSymbolArr, amountArr
}
