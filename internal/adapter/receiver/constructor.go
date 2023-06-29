package receiver

import (
	"sync"

	"github.com/Goboolean/fetch-server/internal/domain/port/in"
)

type StockReceiveAdapter struct {
	port in.RelayerPort
}

var (
	instance *StockReceiveAdapter
	once     sync.Once
)

func New(port in.RelayerPort) *StockReceiveAdapter {
	once.Do(func() {
		instance = &StockReceiveAdapter{
			port: port,
		}
	})

	return instance
}
