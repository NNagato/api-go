package common

import "math/big"

type Token struct {
	Symbol   string `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Decimals uint64 `json:"decimals"`
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
