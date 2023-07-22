package in

import (
	"context"

	"github.com/Goboolean/fetch-server/internal/domain/entity"
)

type ConfiguratorPort interface {
	SetStockRelayableTrue(context.Context, string) error
	SetStockRelayableFalse(context.Context, string) error
	SetStockStoreableTrue(context.Context, string) error
	SetStockStoreableFalse(context.Context, string) error
	SetStockTransmittableTrue(context.Context, string) error
	SetStockTransmittableFalse(context.Context, string) error

	GetStockConfiguration(context.Context, string) (entity.StockConfiguration, error)
	GetAllStockConfiguration(context.Context) ([]entity.StockConfiguration, error)
}
