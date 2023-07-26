package cache

import (
	"context"

	"github.com/Goboolean/fetch-server/internal/domain/entity"
	"github.com/Goboolean/fetch-server/internal/domain/port/out"
)




type MockAdapter struct {
	cache map[string][]*entity.StockAggregate
}

func NewMockAdapter() out.StockPersistenceCachePort {
	return &MockAdapter{
		cache: make(map[string][]*entity.StockAggregate),
	}
}

func (a *MockAdapter) StoreStockOnCache(ctx context.Context, stockId string, stock *entity.StockAggregate) error {

	if _, ok := a.cache[stockId]; !ok {
		a.cache[stockId] = make([]*entity.StockAggregate, 0)		
	}

	a.cache[stockId] = append(a.cache[stockId], stock)

	return nil
}

func (a *MockAdapter) StoreStockBatchOnCache(ctx context.Context, stockId string, batch []*entity.StockAggregate) error {

	if _, ok := a.cache[stockId]; !ok {
		a.cache[stockId] = make([]*entity.StockAggregate, 0)		
	}

	a.cache[stockId] = append(a.cache[stockId], batch...)

	return nil
}

func (a *MockAdapter) GetAndEmptyCache(ctx context.Context, stockId string) ([]*entity.StockAggregate, error) {

	batch := a.cache[stockId]
	a.cache[stockId] = make([]*entity.StockAggregate, 0)

	return batch, nil
}



func (a *MockAdapter) GetStoredStockCount(stockId string) int {
	if _, ok := a.cache[stockId]; !ok {
		a.cache[stockId] = make([]*entity.StockAggregate, 0)		
	}
	return len(a.cache[stockId])
}

func (a *MockAdapter) Clear() {
	for k := range a.cache {
		delete(a.cache, k)
	}
}