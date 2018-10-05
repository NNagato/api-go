package fetcher

import (
	"math/big"

	"github.com/KyberNetwork/api-server/cmd/config"
)

const (
	ETH_TO_WEI = 1000000000000000000
	MIN_ETH    = 0.001
)

func getAmountInWei(amount float64) *big.Int {
	amountFloat := big.NewFloat(amount)
	ethFloat := big.NewFloat(ETH_TO_WEI)
	weiFloat := big.NewFloat(0).Mul(amountFloat, ethFloat)
	amoutInt, _ := weiFloat.Int(nil)
	return amoutInt
}

func MakeDataGetRate(info *config.Config) (
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
	minAmountETH := getAmountInWei(MIN_ETH)

	for _, token := range info.Tokens {
		sourceAddr = append(sourceAddr, token.Address)
		sourceSymbol = append(sourceSymbol, token.Symbol)
		destAddr = append(destAddr, ethAddr)
		destSymbol = append(destSymbol, ethSymbol)
		amount = append(amount, big.NewInt(0))
		amountETH = append(amountETH, minAmountETH)
	}

	sourceArr := append(sourceAddr, destAddr...)
	sourceSymbolArr := append(sourceSymbol, destSymbol...)
	destArr := append(destAddr, sourceAddr...)
	destSymbolArr := append(destSymbol, sourceSymbol...)
	amountArr := append(amount, amountETH...)

	return sourceArr, sourceSymbolArr, destArr, destSymbolArr, amountArr
}
