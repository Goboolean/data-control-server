package relayer

import (
	"context"
	"sync"

	"github.com/Goboolean/fetch-server/internal/domain/entity"
	"github.com/Goboolean/fetch-server/internal/domain/port"
	"github.com/Goboolean/fetch-server/internal/domain/port/out"
)

type RelayerManager struct {
	*store
	*subscriber
	*pipe

	ctx    context.Context
	cancel context.CancelFunc

	tx port.TX
}

var (
	instance *RelayerManager
	once     sync.Once
)

func New(db out.StockPersistencePort, tx port.TX, meta out.StockMetadataPort, ws out.RelayerPort) *RelayerManager {

	once.Do(func() {

		ctx, cancel := context.WithCancel(context.Background())

		startPoint := make(chan *entity.StockAggregateForm)
		endPoint := make(map[string]chan []*entity.StockAggregate)

		instance = &RelayerManager{
			ctx: ctx, cancel: cancel,
			store:      &store{},
			subscriber: newSubscriber(ws, meta, tx),
			pipe:       newPipe(startPoint, endPoint),
		}

		go instance.pipe.ExecPipe(ctx)
	})

	return instance
}

func (m *RelayerManager) Close() error {
	m.cancel()
	return nil
}
