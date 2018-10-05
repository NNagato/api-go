package server

import "github.com/KyberNetwork/api-server/common"

type FetcherInteface interface {
	GetListTokenAPI() ([]common.Token, error)
}
