package relayer

import (
	outport "github.com/Goboolean/stock-fetch-server/internal/domain/port/out"
	"github.com/Goboolean/stock-fetch-server/internal/domain/value"
)

type RawRelayer struct {
	ws outport.RelayerPort

	queue chan value.StockAggregateForm
}

func NewRawRelayer() RawRelayer {
	return RawRelayer{
		queue: make(chan value.StockAggregateForm),
	}
}

func (r *RawRelayer) SubscribeWebsocket(stock string) error {
	return r.ws.SubscribeWebsocket(stock)
}

func (r *RawRelayer) UnsubscribeWebsocket(stock string) error {
	return r.ws.UnsubscribeWebsocket(stock)
}

func (r *RawRelayer) PlaceStockFormBatch(batch []value.StockAggregateForm) error {
	for idx := range batch {
		r.queue <- batch[idx]
	}
	return nil
}
