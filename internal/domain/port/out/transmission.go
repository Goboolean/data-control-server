package out

import (
	"context"

	"github.com/Goboolean/fetch-server/internal/domain/entity"
)

type TransmissionPort interface {
	TransmitStockBatch(ctx context.Context, stockId string, batch []*entity.StockAggregate) error
}
