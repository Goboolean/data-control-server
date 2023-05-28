package config

import (
	"context"
	"fmt"
	"sync"

	api "github.com/Goboolean/stock-fetch-server/api/grpc"
	"github.com/Goboolean/stock-fetch-server/internal/domain/port/in"
	"github.com/Goboolean/stock-fetch-server/internal/domain/service/config"
	"github.com/Goboolean/stock-fetch-server/internal/infrastructure/prometheus"
)

type StockConfiguratorAdapter struct {
	service in.ConfiguratorPort
	api.UnimplementedStockConfiguratorServer
}

var (
	instance *StockConfiguratorAdapter
	once sync.Once
)



func New(s *config.ConfigurationManager) *StockConfiguratorAdapter {

	once.Do(func() {
		instance = &StockConfiguratorAdapter{service: s}
	})

	return instance
}

func (c *StockConfiguratorAdapter) UpdateStockConfiguration(ctx context.Context, in *api.StockConfigUpdateRequest) (nil *api.StockConfigUpdateResponse, err error) {

	prometheus.RequestCounter.Inc()

	stock, optionType, status := in.StockName, in.OptionType, in.OptionStatus

	switch int(optionType) {

	case int(api.StockRelay):
		if status {
			err = c.service.SetStockRelayableTrue(stock)
		} else {
			err = c.service.SetStockRelayableFalse(stock)
		}
		return

	case int(api.StockReal):
		if status {
			err = c.service.SetStockTransmittableTrue(stock)
		} else {
			err = c.service.SetStockTransmittableFalse(stock)
		}
		return

	case int(api.StockPersistance):
		if status {
			err = c.service.SetStockStoreableTrue(stock)
		} else {
			err = c.service.SetStockStoreableFalse(stock)
		}
		return

	default:
		return nil, fmt.Errorf("invalid option type")
	}
}
