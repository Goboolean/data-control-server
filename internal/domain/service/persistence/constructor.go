package persistence

import (
	"context"
	"sync"

	"github.com/Goboolean/fetch-server.v1/internal/domain/port"
	"github.com/Goboolean/fetch-server.v1/internal/domain/port/out"
	"github.com/Goboolean/fetch-server.v1/internal/domain/service/relay"
	relayer "github.com/Goboolean/fetch-server.v1/internal/domain/service/relay"
	"github.com/Goboolean/fetch-server.v1/internal/domain/service/store"
)

type Manager struct {
	o       Option
	db      out.StockPersistencePort
	cache   out.StockPersistenceCachePort
	relayer *relayer.Manager
	s       *store.Store

	tx     port.TX
	ctx    context.Context
	cancel context.CancelFunc
}

var (
	instance *Manager
	once     sync.Once
)

func New(tx port.TX, db out.StockPersistencePort, cache out.StockPersistenceCachePort, r *relay.Manager, o Option) (*Manager, error) {

	ctx, cancel := context.WithCancel(context.Background())

	once.Do(func() {
		instance = &Manager{
			o:       o,
			db:      db,
			cache:   cache,
			relayer: r,
			s:       store.New(ctx),
			tx:      tx,

			ctx:    ctx,
			cancel: cancel,
		}

		if o.BatchSize == 0 {
			o.BatchSize = 1
		}
	})

	instance.setUseCache()
	return instance, nil
}

func (m *Manager) Close() {
	m.cancel()
}
