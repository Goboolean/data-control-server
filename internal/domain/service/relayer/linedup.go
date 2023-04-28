package relayer

import (
	"errors"

	outport "github.com/Goboolean/stock-fetch-server/internal/domain/port/out"
	"github.com/Goboolean/stock-fetch-server/internal/domain/value"
)

type LinedUpRelayer struct {
	ws outport.RelayerPort

	queue map[string]chan []value.StockAggregate
}

func NewLinedUpRelayer() LinedUpRelayer {
	return LinedUpRelayer{
		queue: make(map[string]chan []value.StockAggregate),
	}
}

func (m *LinedUpRelayer) isExist(stock string) bool {
	_, ok := m.queue[stock]
	return ok
}

func (m *LinedUpRelayer) openChannel(stock string) {
	m.queue[stock] = make(chan []value.StockAggregate)
}

func (m *LinedUpRelayer) closeChannel(stock string) {
	delete(m.queue, stock)
}

func (m *LinedUpRelayer) Push(stock string, batch []value.StockAggregate) {
	m.queue[stock] <- batch
}

func (m *LinedUpRelayer) Subscribe(stock string) (<-chan []value.StockAggregate, error) {
	if m.isExist(stock) == false {
		return nil, errors.New("stock does not exist")
	}

	return m.queue[stock], nil
}
