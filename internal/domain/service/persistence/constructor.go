package persistence

import (
	"context"
	"sync"

	"github.com/Goboolean/fetch-server/internal/domain/port"
	"github.com/Goboolean/fetch-server/internal/domain/port/out"
	"github.com/Goboolean/fetch-server/internal/domain/service/relayer"
)

type PersistenceManager struct {
	db      out.StockPersistencePort
	relayer *relayer.RelayerManager
	closed  map[string]chan struct{}

	tx port.TX
	ctx context.Context
	cancel context.CancelFunc
}

var (
	instance *PersistenceManager
	once     sync.Once
)


func New(tx port.TX, db out.StockPersistencePort, r *relayer.RelayerManager) *PersistenceManager {

	once.Do(func() {
		instance = &PersistenceManager{
			db:      db,
			relayer: r,
			closed:  make(map[string]chan struct{}, 1),
			tx:      tx,
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
