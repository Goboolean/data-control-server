package broker

import (
	"context"
	"errors"

	"github.com/Goboolean/fetch-server/internal/domain/entity"
	"github.com/Goboolean/fetch-server/internal/domain/port/out"
)



type MockAdapter struct {
	brokerList map[string] struct{}

}

func NewMockAdapter() out.TransmissionPort {
	return &MockAdapter{

	}
}

func (a MockAdapter) TransmitStockBatch(ctx context.Context, stockId string, batch []*entity.StockAggregate) error {
	if _, ok := a.brokerList[stockId]; !ok {
		return errors.New("path not exist")
	}

	return nil
}

func (a MockAdapter) CreateStockBroker(ctx context.Context, stock string) error {
	a.brokerList[stock] = struct{}{}
	return nil
}
