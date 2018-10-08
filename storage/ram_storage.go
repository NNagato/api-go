package storage

import (
	"sync"

	"github.com/KyberNetwork/api-server/common"
)

type RamStorage struct {
	mu *sync.RWMutex

	rates     []common.Rate
	isNewRate bool
}

func NewRamStorage() *RamStorage {
	mu := &sync.RWMutex{}
	rates := []common.Rate{}
	return &RamStorage{
		mu:        mu,
		rates:     rates,
		isNewRate: false,
	}
}

func (self *RamStorage) SaveRate(rates []common.Rate) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.rates = rates
	self.isNewRate = true
}

func (self *RamStorage) GetRate() []common.Rate {
	self.mu.Lock()
	defer self.mu.Unlock()
	return self.rates
}

func (self *RamStorage) SetIsNewRate(isNewRate bool) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.isNewRate = isNewRate
}

func (self *RamStorage) GetIsNewRate() bool {
	self.mu.Lock()
	defer self.mu.Unlock()
	return self.isNewRate
}
