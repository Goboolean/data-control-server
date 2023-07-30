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
	o 			Option
	db      out.StockPersistencePort
	cache   out.StockPersistenceCachePort
	relayer *relayer.RelayerManager
	s       *store.Store

	tx      port.TX
	ctx     context.Context
	cancel  context.CancelFunc
}

var (
	instance *PersistenceManager
	once     sync.Once
)


func New(tx port.TX, db out.StockPersistencePort, cache out.StockPersistenceCachePort, r *relayer.RelayerManager, o Option) *PersistenceManager {

	ctx, cancel := context.WithCancel(context.Background())

	once.Do(func() {
		instance = &PersistenceManager{
			o:       o,
			db:      db,
			cache:   cache,
			relayer: r,
			s:       store.New(ctx),
			tx:      tx,

			ctx:     ctx,
			cancel:  cancel,
		}

		if o.BatchSize == 0 {
			o.BatchSize = 1
		}
	})

	instance.setUseCache()
	return instance
}


func (m *PersistenceManager) Close() {
	m.cancel()
}
