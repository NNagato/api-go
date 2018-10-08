package common

import "math/big"

type Token struct {
	Symbol   string `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Decimals int    `json:"decimals"`
	Active   bool   `json:"active"`
	Internal bool   `json:"internal"`
}

type ListTokenResponse struct {
	Data    []Token `json:"data"`
	Success bool    `json:"success"`
}

type Rate struct {
	Source  string `json:"source"`
	Dest    string `json:"dest"`
	Rate    string `json:"rate"`
	Minrate string `json:"minRate"`
}

type RateWrapper struct {
	ExpectedRate []*big.Int `json:"expectedRate"`
	SlippageRate []*big.Int `json:"slippageRate"`
}

type ResultRpc struct {
	Result string `json:"result"`
}

type RateBuy struct {
	SrcID  string    `json:"src_id"`
	DstID  string    `json:"dst_id"`
	SrcQty []float64 `json:"src_qty"`
	DstQty []float64 `json:"dst_qty"`
}

type DataGetRate struct {
	SourceArr       []string
	SourceSymbolArr []string
	DstArr          []string
	DstSymbolArr    []string
	AmountArr       []*big.Int
}
