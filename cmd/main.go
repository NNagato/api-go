package main

import (
	"log"
	"runtime"

	"github.com/KyberNetwork/api-server/cmd/config"
	"github.com/KyberNetwork/api-server/fetcher"
	"github.com/KyberNetwork/api-server/server"
	"github.com/KyberNetwork/api-server/storage"
)

func main() {
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)
	//set log for server
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	config, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	fetcher, err := fetcher.NewFetcher(config)
	if err != nil {
		panic(err)
	}

	storage := storage.NewStorage()

	// core := core.NewCore(fetcher, storage)
	server := server.NewServer(fetcher, storage)
	server.Run(config.Port)
}
