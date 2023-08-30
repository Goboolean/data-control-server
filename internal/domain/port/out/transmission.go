package out

import (
	"context"

	"github.com/Goboolean/fetch-server.v1/internal/domain/vo"
)

type TransmissionPort interface {
	TransmitStockBatch(ctx context.Context, stockId string, batch []*vo.StockAggregate) error
	CreateStockBroker(ctx context.Context, stock string) error
}
