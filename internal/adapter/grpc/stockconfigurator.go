package config

import (
	"context"
	"errors"

	"github.com/Goboolean/data-control-server/internal/domain/port/in"
	"github.com/Goboolean/data-control-server/internal/infrastructure/grpc"
)



type StockConfiguratorAdapter struct {
	service inport.StockConfiguratorPort
}

func init() {}



func (c *StockConfiguratorAdapter) UpdateStockConfiguration(ctx context.Context, in *server.StockConfigUpdateRequest) (nil *server.StockConfigUpdateResponse, err error) {
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