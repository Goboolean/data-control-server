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
}

var (
	instance *Transmitter
	once     sync.Once
)

func New(ctx context.Context, broker out.TransmissionPort, relayer *relayer.RelayerManager, o Option) *Transmitter {
	once.Do(func() {
		instance = &Transmitter{
			relayer:   relayer,
			broker:    broker,
			s:         store.New(ctx),
			batchSize: o.BatchSize,
		}
	})

	return instance
}

func (t *Transmitter) Close() {
	t.s.Close()
}
