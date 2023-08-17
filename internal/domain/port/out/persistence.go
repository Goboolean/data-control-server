package out

import (
	"context"

	"github.com/Goboolean/fetch-server/internal/domain/vo"
	"github.com/Goboolean/fetch-server/internal/domain/port"
)

type StockPersistencePort interface {
	StoreStock(port.Transactioner, string, *vo.StockAggregate) error
	StoreStockBatch(port.Transactioner, string, []*vo.StockAggregate) error
	CreateStoringStartedLog(context.Context, string) error
	CreateStoringFailedLog(context.Context, string) error
	CreateStoringStoppedLog(context.Context, string) error
}

type StockPersistenceCachePort interface {
	StoreStockOnCache(context.Context, string, *vo.StockAggregate) error
	StoreStockBatchOnCache(context.Context, string, []*vo.StockAggregate) error
	GetAndEmptyCache(context.Context, string) ([]*vo.StockAggregate, error)
}