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
	defer tx.Rollback()

	if err := m.s.StoreStock(stockId); err != nil {
		return err
	}

	exists, err := m.meta.CheckStockExists(tx, stockId)
	if err != nil {
		return err
	}
	if !exists {
		return ErrStockNotExists
	}

	meta, err := m.meta.GetStockMetadata(tx, stockId)
	if err != nil {
		return err
	}

	if err := m.ws.FetchStock(tx.Context(), stockId, meta); err != nil {
		m.s.UnstoreStock(stockId)
		return err
	}

	return tx.Commit()
}


func (m *RelayerManager) StopFetchingStock(ctx context.Context, stockId string) error {

	if err := m.s.UnstoreStock(stockId); err != nil {
		return err
	}

	if err := m.ws.StopFetchingStock(ctx, stockId); err != nil {
		m.s.StoreStock(stockId)
		return err
	}

	return nil
}


func (m *RelayerManager) IsStockRelayable(stockId string) bool {
	return m.s.StockExists(stockId)
}


func (m *RelayerManager) PlaceStockFormBatch(stockBatch []*entity.StockAggregateForm) {
	for idx := range stockBatch {
		m.pipe.PlaceOnStartPoint(stockBatch[idx])
	}
}


func (m *RelayerManager) Subscribe(stockId string) (<-chan []*entity.StockAggregate, error) {
	return m.pipe.GetEndpointChannel(stockId)
}
