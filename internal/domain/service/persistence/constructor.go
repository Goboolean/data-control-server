package persistence

import (
	"github.com/Goboolean/stock-fetch-server/internal/adapter/stock"
	outport "github.com/Goboolean/stock-fetch-server/internal/domain/port/out"
	"github.com/Goboolean/stock-fetch-server/internal/domain/service/relayer"
)

type PersistenceManager struct {
	db      outport.StockPersistencePort
	relayer *relayer.RelayerManager
	running map[string]chan struct{}
}

var instance *PersistenceManager

func init() {
	instance = &PersistenceManager{
		db:      stock.NewStockAdapter(),
		running: make(map[string]chan struct{}),
	}
}

func NewPersistenceManager() *PersistenceManager {
	return instance
}
