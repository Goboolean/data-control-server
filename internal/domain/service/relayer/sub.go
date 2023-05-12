package relayer

import (
	"fmt"

	"github.com/Goboolean/stock-fetch-server/internal/domain/port/out"
	"github.com/Goboolean/stock-fetch-server/internal/domain/value"
)



type subscriber struct {
	ws out.RelayerPort
	meta out.StockMetadataPort
}


func newSubscriber(ws out.RelayerPort, meta out.StockMetadataPort) *subscriber {
	return &subscriber{ws: ws, meta: meta}
}


func (s *subscriber) fetchStock(stock string) error {

	flag, err := s.meta.StockExists(stock)
	if err != nil {
		return err
	}
	if !flag {
		return fmt.Errorf("stock not exists")
	}

	types, err := s.meta.GetStockType(stock)
	if err != nil {
		return err
	}

	switch types {
	case value.Domestic:
		return s.ws.FetchDomesticStock(stock)
	case value.International:
		return s.ws.FetchInternationalStock(stock)
	default:
		return fmt.Errorf("asdf")
	}
}


func (s *subscriber) unfetchStock(stock string) error {
	return nil
}