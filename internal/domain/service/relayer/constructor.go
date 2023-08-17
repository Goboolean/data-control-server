package relayer

import (
	"context"
	"sync"

	"github.com/Goboolean/fetch-server/internal/domain/vo"
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

func New(db out.StockPersistencePort, tx port.TX, meta out.StockMetadataPort, ws out.RelayerPort) (*RelayerManager, error) {

	once.Do(func() {

		ctx, cancel := context.WithCancel(context.Background())

		instance = &RelayerManager{
			s:      store.New(ctx),
			ws:     ws,
			meta:   meta,
			tx:     tx,

			ctx: 	  ctx,
			cancel: cancel,
		}

		instance.pipe = newPipe()
	  instance.pipe.sinkChan <- &vo.StockAggregateForm{}
		instance.pipe.ExecPipe(ctx)
	})

	return instance, nil
}

func (m *RelayerManager) Close() {
	m.cancel()
}
