package relayer

import (
	"fmt"

	"github.com/Goboolean/stock-fetch-server/internal/domain/port/out"
	"github.com/Goboolean/stock-fetch-server/internal/domain/value"
)

type RawRelayer struct {
	ws out.RelayerPort
	db out.StockPersistencePort
	meta out.StockMetadataPort

	queue chan value.StockAggregateForm
}

func NewRawRelayer() RawRelayer {
	return RawRelayer{
		queue: make(chan value.StockAggregateForm),
	}
}

func (r *RawRelayer) SubscribeWebsocket(stock string) error {

	types, err := r.meta.GetStockType(stock)

	if err != nil {
		return err
	}

	switch types {
	case value.Domestic:
		return r.ws.FetchDomesticStock(stock)
	case value.International:
		return r.ws.FetchInternationalStock(stock)
	default:
		return fmt.Errorf("stock %s do not exist", stock)
	}
}



func (r *RawRelayer) UnsubscribeWebsocket(stock string) error {
	return nil
}



func (r *RawRelayer) PlaceStockFormBatch(batch []value.StockAggregateForm) error {
	for idx := range batch {
		r.queue <- batch[idx]
	}
	return nil
}
