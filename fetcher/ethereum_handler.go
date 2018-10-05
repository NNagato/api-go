package fetcher

import (
	"errors"
	"log"
	"math/big"
	"strings"

	"github.com/KyberNetwork/api-server/common"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

const (
	WRAPPER_ABI = `[{"constant":true,"inputs":[{"name":"x","type":"bytes14"},{"name":"byteInd","type":"uint256"}],"name":"getInt8FromByte","outputs":[{"name":"","type":"int8"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[{"name":"reserve","type":"address"},{"name":"tokens","type":"address[]"}],"name":"getBalances","outputs":[{"name":"","type":"uint256[]"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"pricingContract","type":"address"},{"name":"tokenList","type":"address[]"}],"name":"getTokenIndicies","outputs":[{"name":"","type":"uint256[]"},{"name":"","type":"uint256[]"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"x","type":"bytes14"},{"name":"byteInd","type":"uint256"}],"name":"getByteFromBytes14","outputs":[{"name":"","type":"bytes1"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[{"name":"network","type":"address"},{"name":"sources","type":"address[]"},{"name":"dests","type":"address[]"},{"name":"qty","type":"uint256[]"}],"name":"getExpectedRates","outputs":[{"name":"expectedRate","type":"uint256[]"},{"name":"slippageRate","type":"uint256[]"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"pricingContract","type":"address"},{"name":"tokenList","type":"address[]"}],"name":"getTokenRates","outputs":[{"name":"","type":"uint256[]"},{"name":"","type":"uint256[]"},{"name":"","type":"int8[]"},{"name":"","type":"int8[]"},{"name":"","type":"uint256[]"}],"payable":false,"stateMutability":"view","type":"function"}]`
	NETWORK_ABI = `[{"constant":false,"inputs":[{"name":"alerter","type":"address"}],"name":"removeAlerter","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"reserve","type":"address"},{"name":"src","type":"address"},{"name":"dest","type":"address"},{"name":"add","type":"bool"}],"name":"listPairForReserve","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"},{"name":"","type":"bytes32"}],"name":"perReserveListedPairs","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"getReserves","outputs":[{"name":"","type":"address[]"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"enabled","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"pendingAdmin","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"getOperators","outputs":[{"name":"","type":"address[]"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"token","type":"address"},{"name":"amount","type":"uint256"},{"name":"sendTo","type":"address"}],"name":"withdrawToken","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"maxGasPrice","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"newAlerter","type":"address"}],"name":"addAlerter","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"negligibleRateDiff","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"feeBurnerContract","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"expectedRateContract","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"whiteListContract","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"user","type":"address"}],"name":"getUserCapInWei","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"newAdmin","type":"address"}],"name":"transferAdmin","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"_enable","type":"bool"}],"name":"setEnable","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[],"name":"claimAdmin","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"isReserve","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"getAlerters","outputs":[{"name":"","type":"address[]"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"src","type":"address"},{"name":"dest","type":"address"},{"name":"srcQty","type":"uint256"}],"name":"getExpectedRate","outputs":[{"name":"expectedRate","type":"uint256"},{"name":"slippageRate","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"","type":"uint256"}],"name":"reserves","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"newOperator","type":"address"}],"name":"addOperator","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"reserve","type":"address"},{"name":"add","type":"bool"}],"name":"addReserve","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"operator","type":"address"}],"name":"removeOperator","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"_whiteList","type":"address"},{"name":"_expectedRate","type":"address"},{"name":"_feeBurner","type":"address"},{"name":"_maxGasPrice","type":"uint256"},{"name":"_negligibleRateDiff","type":"uint256"}],"name":"setParams","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"src","type":"address"},{"name":"dest","type":"address"},{"name":"srcQty","type":"uint256"}],"name":"findBestRate","outputs":[{"name":"","type":"uint256"},{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"src","type":"address"},{"name":"srcAmount","type":"uint256"},{"name":"dest","type":"address"},{"name":"destAddress","type":"address"},{"name":"maxDestAmount","type":"uint256"},{"name":"minConversionRate","type":"uint256"},{"name":"walletId","type":"address"}],"name":"trade","outputs":[{"name":"","type":"uint256"}],"payable":true,"stateMutability":"payable","type":"function"},{"constant":false,"inputs":[{"name":"amount","type":"uint256"},{"name":"sendTo","type":"address"}],"name":"withdrawEther","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"getNumReserves","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"token","type":"address"},{"name":"user","type":"address"}],"name":"getBalance","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"admin","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"inputs":[{"name":"_admin","type":"address"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"payable":true,"stateMutability":"payable","type":"fallback"},{"anonymous":false,"inputs":[{"indexed":true,"name":"sender","type":"address"},{"indexed":false,"name":"amount","type":"uint256"}],"name":"EtherReceival","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"source","type":"address"},{"indexed":false,"name":"dest","type":"address"},{"indexed":false,"name":"actualSrcAmount","type":"uint256"},{"indexed":false,"name":"actualDestAmount","type":"uint256"}],"name":"ExecuteTrade","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"reserve","type":"address"},{"indexed":false,"name":"add","type":"bool"}],"name":"AddReserveToNetwork","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"reserve","type":"address"},{"indexed":false,"name":"src","type":"address"},{"indexed":false,"name":"dest","type":"address"},{"indexed":false,"name":"add","type":"bool"}],"name":"ListReservePairs","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"token","type":"address"},{"indexed":false,"name":"amount","type":"uint256"},{"indexed":false,"name":"sendTo","type":"address"}],"name":"TokenWithdraw","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"amount","type":"uint256"},{"indexed":false,"name":"sendTo","type":"address"}],"name":"EtherWithdraw","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"pendingAdmin","type":"address"}],"name":"TransferAdminPending","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"newAdmin","type":"address"},{"indexed":false,"name":"previousAdmin","type":"address"}],"name":"AdminClaimed","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"newAlerter","type":"address"},{"indexed":false,"name":"isAdd","type":"bool"}],"name":"AlerterAdded","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"newOperator","type":"address"},{"indexed":false,"name":"isAdd","type":"bool"}],"name":"OperatorAdded","type":"event"}]`
)

type EthereumHandler struct {
	networkAddr ethereum.Address
	networkAbi  abi.ABI
	wrapperAddr ethereum.Address
	wrapperAbi  abi.ABI
}

func NewEthereumHandler(networkAddr, wrapperAddr string) (*EthereumHandler, error) {
	wrapperAbi, err := abi.JSON(strings.NewReader(WRAPPER_ABI))
	if err != nil {
		log.Print(err)
		return nil, err
	}
	networkAbi, err := abi.JSON(strings.NewReader(NETWORK_ABI))
	if err != nil {
		log.Print(err)
		return nil, err
	}

	ethereum := &EthereumHandler{
		networkAddr: ethereum.HexToAddress(networkAddr),
		networkAbi:  networkAbi,
		wrapperAddr: ethereum.HexToAddress(wrapperAddr),
		wrapperAbi:  wrapperAbi,
	}
	return ethereum, nil
}

// --------------------- Encode Data ---------------------------------------

func (self *EthereumHandler) EncodeRateData(sourceArr, destArr []string, quantity []*big.Int) ([]byte, error) {
	sourceList := make([]ethereum.Address, 0)
	for _, sourceItem := range sourceArr {
		sourceList = append(sourceList, ethereum.HexToAddress(sourceItem))
	}
	destList := make([]ethereum.Address, 0)
	for _, destItem := range destArr {
		destList = append(destList, ethereum.HexToAddress(destItem))
	}

	quantityList := make([]*big.Int, 0)
	for _, quanItem := range quantity {
		quantityList = append(quantityList, quanItem)
	}

	encodedData, err := self.wrapperAbi.Pack("getExpectedRates", self.networkAddr, sourceList, destList, quantityList)
	if err != nil {
		log.Print(err)
		return []byte{}, err
	}

	return encodedData, nil
}

// ---------------------- extract data ---------------------------------------

func (self *EthereumHandler) ExtractRateData(result string, sourceArr []string, destAddr []string) (*[]common.Rate, error) {
	rateByte, err := hexutil.Decode(result)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	var rateWapper common.RateWrapper
	err = self.wrapperAbi.Unpack(&rateWapper, "getExpectedRates", rateByte)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	var lenArr = len(sourceArr)
	if (len(rateWapper.ExpectedRate) != lenArr) || (len(rateWapper.SlippageRate) != lenArr) {
		errorLength := errors.New("Length of expected for slippage rate is not enough")
		log.Print(errorLength)
		return nil, errorLength
	}

	rateReturn := make([]common.Rate, 0)
	for i := 0; i < lenArr; i++ {
		source := sourceArr[i]
		dest := destAddr[i]
		rate := rateWapper.ExpectedRate[i]
		minRate := rateWapper.SlippageRate[i]
		rateReturn = append(rateReturn, common.Rate{
			source, dest, rate.String(), minRate.String(),
		})
	}
	return &rateReturn, nil
}
