package main

import (
	"log"
	"runtime"
	"time"

	core "github.com/KyberNetwork/api-server/api-core"
	"github.com/KyberNetwork/api-server/cmd/config"
	"github.com/KyberNetwork/api-server/fetcher"
	"github.com/KyberNetwork/api-server/server"
	"github.com/KyberNetwork/api-server/storage"
)

type fetcherFunc func(*storage.RamStorage, *fetcher.Fetcher)

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

	ramStorage := storage.NewRamStorage()
	core := core.NewCore(fetcher, ramStorage)
	server := server.NewServer(core)

	// run fetcher
	runFetchData(ramStorage, fetcher, FetchRate, 20)

	server.Run(config.Port)
}

func runFetchData(storage *storage.RamStorage, fertcherIns *fetcher.Fetcher, fn fetcherFunc, interval time.Duration) {
	ticker := time.NewTicker(interval * time.Second)
	go func() {
		for {
			fn(storage, fertcherIns)
			<-ticker.C
		}
	}()
}

func FetchRate(storage *storage.RamStorage, fertcher *fetcher.Fetcher) {
	currentRate := storage.GetRate()
	rates, err := fertcher.GetRate(currentRate, nil)
	if err != nil {
		log.Println(err)
		storage.SetIsNewRate(false)
		return
	}
	storage.SaveRate(*rates)
}
