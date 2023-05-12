package relayer

import (
	"context"
	"sync"

	"github.com/Goboolean/stock-fetch-server/internal/domain/port/out"
	"github.com/Goboolean/stock-fetch-server/internal/domain/value"
)


type RelayerManager struct {
	*store
	*subscriber
	*pipe

	ctx context.Context
	cancel context.CancelFunc
}



var (
	instance *RelayerManager
	once sync.Once
)

func New(db out.StockPersistencePort, meta out.StockMetadataPort, ws out.RelayerPort) *RelayerManager {

	once.Do(func() {

		ctx, cancel := context.WithCancel(context.Background())

		startPoint := make(chan value.StockAggregateForm)
		endPoint := make(map[string] chan []value.StockAggregate)

		instance = &RelayerManager{
			ctx: ctx, cancel: cancel,
			store: &store{},
			subscriber: newSubscriber(ws, meta),
			pipe: newPipe(startPoint, endPoint),
		}

		go instance.pipe.ExecPipe(ctx)
	})

	return instance
}



func (m *RelayerManager) Close() error {
	m.cancel()
	return nil
}

