package persistence

import (
	"sync"

	"github.com/Goboolean/stock-fetch-server/internal/domain/port/out"
	"github.com/Goboolean/stock-fetch-server/internal/domain/service/relayer"
)

type PersistenceManager struct {
	db      out.StockPersistencePort
	relayer *relayer.RelayerManager
	closed map[string]chan struct{}
}

var (
	instance *PersistenceManager
	once sync.Once
)



func New(db out.StockPersistencePort, r *relayer.RelayerManager) *PersistenceManager {

	once.Do(func() {
		instance = &PersistenceManager{
			db:      db,
			relayer: r,
			closed: make(map[string]chan struct{}, 1),
		}
	})

	return instance
}



func (m *PersistenceManager) Close() {
	for stock, ch := range m.closed {
		ch <- struct{}{}
		delete(m.closed, stock)
	}
}