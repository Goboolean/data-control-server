package relay

import (
	"context"
	"sync"

	"github.com/Goboolean/fetch-server.v1/internal/domain/port"
	"github.com/Goboolean/fetch-server.v1/internal/domain/port/out"
	"github.com/Goboolean/fetch-server.v1/internal/domain/service/store"
	"github.com/Goboolean/fetch-server.v1/internal/domain/vo"
)

type Manager struct {
	s    *store.Store
	ws   out.RelayerPort
	meta out.StockMetadataPort

	pipe *pipe

	ctx    context.Context
	cancel context.CancelFunc

	tx port.TX
}

var (
	instance *Manager
	once     sync.Once
)

func New(db out.StockPersistencePort, tx port.TX, meta out.StockMetadataPort, ws out.RelayerPort) (*Manager, error) {

	once.Do(func() {

		ctx, cancel := context.WithCancel(context.Background())

		instance = &Manager{
			s:    store.New(ctx),
			ws:   ws,
			meta: meta,
			tx:   tx,

			ctx:    ctx,
			cancel: cancel,
		}

		instance.pipe = newPipe()
		instance.pipe.sinkChan <- &vo.StockAggregateForm{}
		instance.pipe.ExecPipe(ctx)
	})

	return instance, nil
}

func (m *Manager) Close() {
	m.cancel()
}
