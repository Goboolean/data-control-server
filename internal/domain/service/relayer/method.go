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

	exists, err := m.meta.CheckStockExists(tx, stockId)
	if err != nil {
		return err
	}
	if !exists {
		return ErrStockNotExists
	}

	if err := m.s.StoreStock(stockId); err != nil {
		return err
	}

	m.pipe.AddNewPipe(stockId)

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

	m.pipe.RemovePipe(stockId)

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

	if exists := m.s.StockExists(stockId); !exists {
		return nil, ErrStockNotExists
	}

	return m.pipe.RegisterNewSubscriber(ctx, stockId)
}
