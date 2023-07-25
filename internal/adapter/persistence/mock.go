package persistence

import (
	"context"

	"github.com/Goboolean/fetch-server/internal/domain/entity"
	"github.com/Goboolean/fetch-server/internal/domain/port"
	"github.com/Goboolean/fetch-server/internal/domain/port/out"
)




type MockAdapter struct {}

func NewMockAdapter() out.StockPersistencePort {
	return &MockAdapter{}
}


func (a *MockAdapter) StoreStock(port.Transactioner, string, *entity.StockAggregate) error {
	return nil
}

func (a *MockAdapter) StoreStockBatch(port.Transactioner, string, []*entity.StockAggregate) error {
	return nil
}

func (a *MockAdapter) CreateStoringStartedLog(context.Context, string) error {
	return nil
}

func (a *MockAdapter) CreateStoringFailedLog(context.Context, string) error {
	return nil
}

func (a *MockAdapter) CreateStoringStoppedLog(context.Context, string) error {
	return nil
}