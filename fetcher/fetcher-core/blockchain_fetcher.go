package fetchercore

import (
	"context"
	"log"
	"math/big"

	ether "github.com/ethereum/go-ethereum"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type BlockchainFetcher struct {
	client      *rpc.Client
	ethclient   *ethclient.Client
	networkAddr ethereum.Address
	wrapperAddr ethereum.Address
}

func NewBlockchainFetcher(endpoint, networkAddr, wrapperAddr string) (*BlockchainFetcher, error) {
	client, err := rpc.Dial(endpoint)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	ethclient := ethclient.NewClient(client)
	return &BlockchainFetcher{
		client:      client,
		ethclient:   ethclient,
		networkAddr: ethereum.HexToAddress(networkAddr),
		wrapperAddr: ethereum.HexToAddress(wrapperAddr),
	}, nil
}

func (self *BlockchainFetcher) EthCall(toAddr string, dataABI string) (string, error) {
	params := make(map[string]string)
	params["data"] = dataABI
	params["to"] = toAddr

	var result string
	err := self.client.Call(&result, "eth_call", params, "latest")
	if err != nil {
		log.Print(err)
		return "", err
	}

	return result, nil
}

// ------------- other way to GetRate
// func (self *BlockchainFetcher) GetRate(dataABI []byte) (string, error) {
// 	msg := ether.CallMsg{
// 		To:   &self.wrapperAddr,
// 		Data: dataABI,
// 	}
// 	ctx := context.Background()
// 	rateByte, err := self.ethclient.CallContract(ctx, msg, nil)
// 	if err != nil {
// 		log.Println(err)
// 		return "0x0", err
// 	}
// 	return ethereum.ToHex(rateByte), nil
// }

// ----------- Get balance in WEI

func (self *BlockchainFetcher) GetBalanceAccount() {
	balance, err := self.ethclient.BalanceAt(context.Background(), self.wrapperAddr, nil)
	log.Println("balance account: ", balance, err)
}

// ----------- EstimateGas example

func (self *BlockchainFetcher) EstimateGas() {
	from := "0x2cc72d8857ac57ba058eddd36b2f14adc2"
	data := "0xcb3c28c7000000000000000000000000eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee00000000000000000000000000000000000000000000000000038d7ea4c680000000000000000000000000004e470dc7321e84ca96fcaedd0c8abcebbaeb68c60000000000000000000000002cc72d8857ac57ba058eddd36b2f14adc2a058bd800000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001431b0294fa51ef00000000000000000000000000000000000000000000000000000000000003e1e4b"
	to := "0x9E67f627a17EDed3Fb7C71417DFe4aa7bFb4CaB7"
	toAddress := ethereum.HexToAddress(to)
	value := big.NewInt(1000000000)
	dataByte, err := hexutil.Decode(data)
	if err != nil {
		log.Println("cant decode data: ", err)
	}

	msg := ether.CallMsg{
		From:  ethereum.HexToAddress(from),
		To:    &toAddress,
		Value: value,
		Data:  dataByte,
	}

	result, err := self.ethclient.EstimateGas(context.Background(), msg)
	log.Println("Estimate gas result: ", result, err)
}
