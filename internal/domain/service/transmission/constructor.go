package transmission

import (
	"sync"

	"github.com/Goboolean/stock-fetch-server/internal/domain/port/out"
	"github.com/Goboolean/stock-fetch-server/internal/domain/service/relayer"
	"github.com/Goboolean/stock-fetch-server/internal/domain/value"
)


type Transmitter struct {
	relayer *relayer.RelayerManager
	broker out.TransmissionPort

	closed map[string] chan []value.StockAggregate
}

var (
	instance *Transmitter
	once sync.Once
)

func New(broker out.TransmissionPort, relayer *relayer.RelayerManager) *Transmitter {
	once.Do(func() {
		instance = &Transmitter{
			relayer: relayer,
			broker: broker,
		}
	})

	return instance
}