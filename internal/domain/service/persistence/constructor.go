package persistence

import (
	"sync"

	outport "github.com/Goboolean/stock-fetch-server/internal/domain/port/out"
	"github.com/Goboolean/stock-fetch-server/internal/domain/service/relayer"
)

type PersistenceManager struct {
	db      outport.StockPersistencePort
	relayer *relayer.RelayerManager
	running map[string]chan struct{}
}

var (
	instance *PersistenceManager
	once sync.Once
)



func New(db outport.StockPersistencePort, relayer *relayer.RelayerManager) *PersistenceManager {

	once.Do(func() {
		instance = &PersistenceManager{
			db:      db,
			relayer: relayer,
			running: make(map[string]chan struct{}),
		}
	})

	return instance
}
