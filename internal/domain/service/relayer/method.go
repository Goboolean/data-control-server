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

	m.pipe.AddNewPipe(stockId)

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

	if err := m.ws.FetchStock(tx.Context(), meta.StockID, meta.Platform, meta.Symbol); err != nil {
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
	for _, stock := range stockBatch {
		m.pipe.PlaceOnStartPoint(stock)
	}
}

// 
// If call side execute ctx.Done(), then subscription of this stock will be cancelled.
func (m *RelayerManager) Subscribe(ctx context.Context, stockId string) (<-chan *entity.StockAggregate, error) {
	return m.pipe.RegisterNewSubscriber(ctx, stockId)
}
