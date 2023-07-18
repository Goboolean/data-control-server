package relayer

import (
	"context"

	"github.com/Goboolean/fetch-server/internal/domain/entity"
)


func (m *RelayerManager) FetchStock(ctx context.Context, stockId string) error {

	tx, err := m.tx.Transaction(ctx)
	if err != nil {
		return err
	}

	if err := m.store.storeStock(stockId); err != nil {
		return err
	}

	if err := m.subscriber.fetchStock(tx, stockId); err != nil {
		m.store.unstoreStock(stockId)
		return err
	}

	return nil
}


func (m *RelayerManager) StopFetchingStock(stockId string) error {
	if err := m.store.unstoreStock(stockId); err != nil {
		return err
	}

	if err := m.subscriber.unfetchStock(stockId); err != nil {
		m.store.storeStock(stockId)
		return err
	}

	return nil
}


func (m *RelayerManager) PlaceStockFormBatch(stockBatch []*entity.StockAggregateForm) {
	for idx := range stockBatch {
		m.pipe.PlaceOnStartPoint(stockBatch[idx])
	}
}


func (m *RelayerManager) Subscribe(stockId string) (<-chan []*entity.StockAggregate, error) {
	return m.pipe.GetEndpointChannel(stockId)
}
