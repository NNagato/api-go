package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/KyberNetwork/api-server/common"
)

type Config struct {
	NodeEndpoint   []string       `json:"nodeEndpoint"`
	EtherScanAPI   string         `json:"etherScanAPI"`
	ETHAddr        string         `json:"eth_address"`
	ETHSymbol      string         `json:"eth_symbol"`
	CMCendpoint    string         `json:"api_usd"`
	Reserver       string         `json:"reserve"`
	Network        string         `json:"network"`
	Wrapper        string         `json:"wrapper"`
	TrackerEnpoint string         `json:"tracker_endpoint"`
	CoreEndpoint   string         `json:"core_endpoint"`
	Tokens         []common.Token `json:"tokens"`

	SecretCore string `json:"secret_core"`

	Port string `json:"port"`
}

func NewConfig() (*Config, error) {
	var file []byte
	var err error
	switch os.Getenv("KYBER_ENV") {
	case "internal_mainnet":
		file, err = ioutil.ReadFile("../env/internal_mainnet.json")
		if err != nil {
			log.Print(err)
			return nil, err
		}
		break
	case "staging":
		file, err = ioutil.ReadFile("../env/staging.json")
		if err != nil {
			log.Print(err)
			return nil, err
		}
		break
	case "production":
		file, err = ioutil.ReadFile("../env/production.json")
		if err != nil {
			log.Print(err)
			return nil, err
		}
	case "kovan":
		file, err = ioutil.ReadFile("../env/kovan.json")
		if err != nil {
			log.Print(err)
			return nil, err
		}
	case "ropsten":
		file, err = ioutil.ReadFile("../env/ropsten.json")
		if err != nil {
			log.Print(err)
			return nil, err
		}
	case "production_test":
		file, err = ioutil.ReadFile("../env/production_test.json")
		if err != nil {
			log.Print(err)
			return nil, err
		}
	default:
		file, err = ioutil.ReadFile("../env/production.json")
		if err != nil {
			log.Print(err)
			return nil, err
		}
		break
	}

	config := &Config{}
	err = json.Unmarshal(file, config)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return config, nil
}
