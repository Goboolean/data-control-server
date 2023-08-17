package broker

import (
	"context"
	"errors"

	"github.com/Goboolean/fetch-server/internal/domain/vo"
	"github.com/Goboolean/fetch-server/internal/domain/port/out"
)



type MockAdapter struct {
	brokerList map[string] int
}

func NewMockAdapter() out.TransmissionPort {
	return &MockAdapter{
		brokerList: make(map[string]int),
	}
}

func (a *MockAdapter) TransmitStockBatch(ctx context.Context, stockId string, batch []*vo.StockAggregate) error {
	if _, ok := a.brokerList[stockId]; !ok {
		return errors.New("path not exist")
	}

	a.brokerList[stockId] += len(batch)
	return nil
}

func (a *MockAdapter) CreateStockBroker(ctx context.Context, stock string) error {
	a.brokerList[stock] = 0
	return nil
}

func (a *MockAdapter) GetTransmittedStockCount(stockId string) (int, error) {
	count, ok := a.brokerList[stockId]
	if !ok {
		return 0, errors.New("stock not exist")
	}
	return count, nil
}