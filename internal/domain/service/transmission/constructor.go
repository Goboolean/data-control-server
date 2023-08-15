package transmission

import (
	"context"
	"sync"

	"github.com/Goboolean/fetch-server/internal/domain/port/out"
	"github.com/Goboolean/fetch-server/internal/domain/service/relayer"
	"github.com/Goboolean/fetch-server/internal/domain/service/store"
)

type Transmitter struct {
	relayer *relayer.RelayerManager
	broker out.TransmissionPort

	s *store.Store
	batchSize int

	ctx context.Context
	cancel context.CancelFunc
}

var (
	instance *Transmitter
	once     sync.Once
)

func New(broker out.TransmissionPort, relayer *relayer.RelayerManager, o Option) (*Transmitter, error) {

	ctx, cancel := context.WithCancel(context.Background())

	once.Do(func() {
		instance = &Transmitter{
			relayer:   relayer,
			broker:    broker,
			s:         store.New(ctx),
			batchSize: o.BatchSize,

			ctx: ctx,
			cancel: cancel,
		}
	})

	return instance, nil
}

func (t *Transmitter) Close() {
	t.cancel()
}
