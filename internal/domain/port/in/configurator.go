package in

import (
	"context"

	"github.com/Goboolean/fetch-server/internal/domain/vo"
)

type ConfiguratorPort interface {
	SetStockRelayableTrue(context.Context, string) error
	SetStockRelayableFalse(context.Context, string) error
	SetStockStoreableTrue(context.Context, string) error
	SetStockStoreableFalse(context.Context, string) error
	SetStockTransmittableTrue(context.Context, string) error
	SetStockTransmittableFalse(context.Context, string) error

	GetStockConfiguration(context.Context, string) (vo.StockConfiguration, error)
	GetAllStockConfiguration(context.Context) ([]vo.StockConfiguration, error)
}
