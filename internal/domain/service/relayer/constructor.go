package relayer

import (
	"context"
	"sync"

	"github.com/Goboolean/fetch-server/internal/domain/entity"
	"github.com/Goboolean/fetch-server/internal/domain/port"
	"github.com/Goboolean/fetch-server/internal/domain/port/out"
	"github.com/Goboolean/fetch-server/internal/domain/service/store"
)

type RelayerManager struct {
	s *store.Store
	ws   out.RelayerPort
	meta out.StockMetadataPort

	pipe *pipe

	ctx    context.Context
	cancel context.CancelFunc

	tx port.TX
}

var (
	instance *RelayerManager
	once     sync.Once
)

func New(ctx context.Context, db out.StockPersistencePort, tx port.TX, meta out.StockMetadataPort, ws out.RelayerPort) *RelayerManager {

	once.Do(func() {

		ctx, cancel := context.WithCancel(ctx)

		instance = &RelayerManager{
			ctx:    ctx,
			cancel: cancel,
			s:      store.New(ctx),
			ws:     ws,
			meta:   meta,
			tx:     tx,
		}

		instance.pipe = newPipe()
	  instance.pipe.sinkChan <- &entity.StockAggregateForm{}
		instance.pipe.ExecPipe(ctx)
	})

	return instance
}

func (m *RelayerManager) Close() error {
	m.cancel()
	return nil
}
