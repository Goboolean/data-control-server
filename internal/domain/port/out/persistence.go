package out

import (
	"context"

	"github.com/Goboolean/fetch-server/internal/domain/entity"
	"github.com/Goboolean/fetch-server/internal/domain/port"
)

type StockPersistencePort interface {
	StoreStock(port.Transactioner, string, *entity.StockAggregate) error
	StoreStockBatch(port.Transactioner, string, []*entity.StockAggregate) error
	CreateStoringStartedLog(context.Context, string) error
	CreateStoringFailedLog(context.Context, string) error
	CreateStoringStoppedLog(context.Context, string) error
}
