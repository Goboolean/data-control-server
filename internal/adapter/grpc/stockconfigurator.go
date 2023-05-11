package config

import (
	"context"
	"errors"

	"github.com/Goboolean/stock-fetch-server/internal/domain/port/in"
	"github.com/Goboolean/stock-fetch-server/internal/domain/service/config"
	api "github.com/Goboolean/stock-fetch-server/api/grpc"
)

type StockConfiguratorAdapter struct {
	service in.ConfiguratorPort
	api.UnimplementedStockConfiguratorServer
}

var instance *StockConfiguratorAdapter

func init() {
	instance = &StockConfiguratorAdapter{service: config.NewConfigurationManager()}
}

func New() *StockConfiguratorAdapter {
	return instance
}

func (c *StockConfiguratorAdapter) UpdateStockConfiguration(ctx context.Context, in *api.StockConfigUpdateRequest) (nil *api.StockConfigUpdateResponse, err error) {
	stock := in.StockName
	optionType := in.OptionType
	status := in.OptionStatus

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
