package config

import (
	"context"
	"errors"
	"sync"

	api "github.com/Goboolean/stock-fetch-server/api/grpc"
	"github.com/Goboolean/stock-fetch-server/internal/domain/port/in"
	"github.com/Goboolean/stock-fetch-server/internal/domain/service/config"
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
			err = c.service.SetStockRelayableTrue(stock)
		} else {
			err = c.service.SetStockRelayableFalse(stock)
		}
		return

	case int(api.StockPersistance):
		if status {
			err = c.service.SetStockRelayableTrue(stock)
		} else {
			err = c.service.SetStockRelayableFalse(stock)
		}
		return

	default:
		return nil, errors.New("invalid option type")
	}
}
