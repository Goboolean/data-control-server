package transmission

import (
	"context"
	"sync"

	"github.com/Goboolean/fetch-server/internal/domain/port/out"
	"github.com/Goboolean/fetch-server/internal/domain/service/relay"
	"github.com/Goboolean/fetch-server/internal/domain/service/store"
)

type Manager struct {
	relayer *relay.Manager
	broker  out.TransmissionPort

	s         *store.Store
	batchSize int

	ctx    context.Context
	cancel context.CancelFunc
}

var (
	instance *Manager
	once     sync.Once
)

func New(broker out.TransmissionPort, relayer *relay.Manager, o Option) (*Manager, error) {

	ctx, cancel := context.WithCancel(context.Background())

	once.Do(func() {
		instance = &Manager{
			relayer:   relayer,
			broker:    broker,
			s:         store.New(ctx),
			batchSize: o.BatchSize,

			ctx:    ctx,
			cancel: cancel,
		}
	})

	return instance, nil
}

func (t *Manager) Close() {
	t.cancel()
}
