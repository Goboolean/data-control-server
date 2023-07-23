package persistence

import (
	"context"
	"sync"

	"github.com/Goboolean/fetch-server/internal/domain/port"
	"github.com/Goboolean/fetch-server/internal/domain/port/out"
	"github.com/Goboolean/fetch-server/internal/domain/service/relayer"
	"github.com/Goboolean/fetch-server/internal/domain/service/store"
)

type PersistenceManager struct {
	db      out.StockPersistencePort
	cache   out.StockPersistenceCachePort
	relayer *relayer.RelayerManager
	s       *store.Store

	tx port.TX
}

var (
	instance *PersistenceManager
	once     sync.Once
)


func New(tx port.TX, ctx context.Context, db out.StockPersistencePort, cache out.StockPersistenceCachePort, r *relayer.RelayerManager) *PersistenceManager {

	once.Do(func() {
		instance = &PersistenceManager{
			db:      db,
			cache:   cache,
			relayer: r,
			s:       store.New(ctx),
			tx:      tx,
		}
	})

	return instance
}


func (m *PersistenceManager) Close() {
	m.s.Close()	
}
