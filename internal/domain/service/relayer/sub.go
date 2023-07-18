package relayer

import (
	"github.com/Goboolean/fetch-server/internal/domain/port"
	"github.com/Goboolean/fetch-server/internal/domain/port/out"
)

type subscriber struct {
	ws   out.RelayerPort
	meta out.StockMetadataPort
	tx port.TX
}


func newSubscriber(ws out.RelayerPort, meta out.StockMetadataPort, tx port.TX) *subscriber {
	return &subscriber{ws: ws, meta: meta}
}


func (s *subscriber) fetchStock(tx port.Transactioner, stockId string) error {

	exists, err := s.meta.CheckStockExists(tx, stockId)
	if err != nil {
		return err
	}
	if !exists {
		return ErrStockNotExists
	}

	meta, err := s.meta.GetStockMetadata(tx, stockId)
	if err != nil {
		return err
	}

	return s.ws.FetchStock(stockId, meta)
}


func (s *subscriber) unfetchStock(stock string) error {
	return s.ws.StopFetchingStock(stock)
}
