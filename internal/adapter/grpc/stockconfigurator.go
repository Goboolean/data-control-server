package config

import (
	"context"
	"errors"

	"github.com/Goboolean/stock-fetch-server/internal/domain/port/in"
	"github.com/Goboolean/stock-fetch-server/internal/domain/service/config"
	model "github.com/Goboolean/stock-fetch-server/internal/infrastructure/grpc/config"
)

type StockConfiguratorAdapter struct {
	service inport.ConfiguratorPort
	model.UnimplementedStockConfiguratorServer
}

var instance *StockConfiguratorAdapter

func init() {
	instance = &StockConfiguratorAdapter{service: config.NewConfigurationManager()}
}

func New() *StockConfiguratorAdapter {
	return instance
}

func (c *StockConfiguratorAdapter) UpdateStockConfiguration(ctx context.Context, in *model.StockConfigUpdateRequest) (nil *model.StockConfigUpdateResponse, err error) {
	stock := in.StockName
	optionType := in.OptionType
	status := in.OptionStatus

	switch int(optionType) {

	case int(StockRelay):
		if status {
			err = c.service.SetStockRelayableTrue(stock)
		} else {
			err = c.service.SetStockRelayableFalse(stock)
		}
		return

	case int(StockReal):
		if status {
			err = c.service.SetStockRelayableTrue(stock)
		} else {
			err = c.service.SetStockRelayableFalse(stock)
		}
		return

	case int(StockPersistance):
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
