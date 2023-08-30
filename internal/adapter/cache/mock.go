package cache

import (
	"context"

	"github.com/Goboolean/fetch-server.v1/internal/domain/port/out"
	"github.com/Goboolean/fetch-server.v1/internal/domain/vo"
)

type MockAdapter struct {
	cache map[string][]*vo.StockAggregate
}

func NewMockAdapter() out.StockPersistenceCachePort {
	return &MockAdapter{
		cache: make(map[string][]*vo.StockAggregate),
	}
}

func (a *MockAdapter) StoreStockOnCache(ctx context.Context, stockId string, stock *vo.StockAggregate) error {

	if _, ok := a.cache[stockId]; !ok {
		a.cache[stockId] = make([]*vo.StockAggregate, 0)
	}

	a.cache[stockId] = append(a.cache[stockId], stock)

	return nil
}

func (a *MockAdapter) StoreStockBatchOnCache(ctx context.Context, stockId string, batch []*vo.StockAggregate) error {

	if _, ok := a.cache[stockId]; !ok {
		a.cache[stockId] = make([]*vo.StockAggregate, 0)
	}

	a.cache[stockId] = append(a.cache[stockId], batch...)

	return nil
}

func (a *MockAdapter) GetAndEmptyCache(ctx context.Context, stockId string) ([]*vo.StockAggregate, error) {

	batch := a.cache[stockId]
	a.cache[stockId] = make([]*vo.StockAggregate, 0)

	return batch, nil
}

func (a *MockAdapter) GetStoredStockCount(stockId string) int {
	if _, ok := a.cache[stockId]; !ok {
		a.cache[stockId] = make([]*vo.StockAggregate, 0)
	}
	return len(a.cache[stockId])
}

func (a *MockAdapter) Clear() {
	for k := range a.cache {
		delete(a.cache, k)
	}
}
