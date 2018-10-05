package fetcher

type FetcherInterfce interface {
	EthCall(toAddr string, dataAbi string) (string, error)
	GetBalanceAccount()
	EstimateGas()
}
