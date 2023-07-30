package grpc



import (
	"context"
	"errors"
	"fmt"

	api "github.com/Goboolean/fetch-server/api/grpc"
	"github.com/Goboolean/fetch-server/internal/domain/port/in"
)



type Adapter struct {
	service in.ConfiguratorPort
	api.UnimplementedStockConfiguratorServer
}

func NewAdapter(s in.ConfiguratorPort) api.StockConfiguratorServer {
	return &Adapter{service: s}
}



func (c *Adapter) UpdateStockConfigOne(ctx context.Context, in *api.StockConfig) (*api.ReturnMessage, error) {

	stockId := in.GetStockId()
	relayable := in.GetRelayable()
	transmittable := in.GetTransmittable()
	storeable := in.GetStoreable()

	result := &api.ReturnMessage{}

	switch relayable.OptionStatus {
	case 1:
		if err := c.service.SetStockRelayableTrue(ctx, stockId); err != nil {
			result.Message = err.Error()
			return result, err
		}
	case 0:
		if err := c.service.SetStockRelayableFalse(ctx, stockId); err != nil {
			result.Message = err.Error()
			return result, err
		}
	case -1:
		break
	default:
		return nil, fmt.Errorf("invalid option status")
	}

	switch transmittable.OptionStatus {
	case 1:
		if err := c.service.SetStockTransmittableTrue(ctx, stockId); err != nil {
			result.Message = err.Error()
			return result, err
		}
	case 0:
		if err := c.service.SetStockTransmittableFalse(ctx, stockId); err != nil {
			result.Message = err.Error()
			return result, err
		}
	case -1:
		break
	default:
		return nil, fmt.Errorf("invalid option status")
	}

	switch storeable.OptionStatus {
	case 1:
		if err := c.service.SetStockStoreableTrue(ctx, stockId); err != nil {
			result.Message = err.Error()
			return result, err
		}
	case 0:
		if err := c.service.SetStockStoreableFalse(ctx, stockId); err != nil {
			result.Message = err.Error()
			return result, err
		}
	case -1:
		break
	default:
		return nil, fmt.Errorf("invalid option status")
	}

	return result, nil
}



func (c *Adapter) UpdateStockConfigMany(ctx context.Context, in *api.StockConfigList) (*api.ReturnMessageList, error) {

	length := len(in.GetStockConfig())

	msgList := make([]*api.ReturnMessage, length)

	for i, v := range in.GetStockConfig() {

		var result *api.ReturnMessage
		
		stockId := v.GetStockId()
		relayable := v.GetRelayable()
		transmittable := v.GetTransmittable()
		storeable := v.GetStoreable()


		switch relayable.OptionStatus {
		case 1:
			if err := c.service.SetStockRelayableTrue(ctx, stockId); err != nil {
				result.Message = err.Error()
				break
			}
		case 0:
			if err := c.service.SetStockRelayableFalse(ctx, stockId); err != nil {
				result.Message = err.Error()
				break
			}
		case -1:
			break
		default:
			return nil, fmt.Errorf("invalid option status")
		}

		switch transmittable.OptionStatus {
		case 1:
			if err := c.service.SetStockTransmittableTrue(ctx, stockId); err != nil {
				result.Message = err.Error()
				break
			}
		case 0:
			if err := c.service.SetStockTransmittableFalse(ctx, stockId); err != nil {
				result.Message = err.Error()
				break
			}
		case -1:
			break
		default:
			return nil, fmt.Errorf("invalid option status")
		}

		switch storeable.OptionStatus {
		case 1:
			if err := c.service.SetStockStoreableTrue(ctx, stockId); err != nil {
				result.Message = err.Error()
				break
			}
		case 0:
			if err := c.service.SetStockStoreableFalse(ctx, stockId); err != nil {
				result.Message = err.Error()
				break
			}
		case -1:
			break
		default:
			return nil, fmt.Errorf("invalid option status")
		}

		msgList[i] = result
	}

	var err error
	for _, v := range msgList {
		if v.Status == false {
			err = errors.New("")
			break
		}
	}

	return &api.ReturnMessageList{ReturnMessage: msgList}, err
}



func (c *Adapter) GetStockConfigOne(ctx context.Context, in *api.StockId) (*api.StockConfig, error) {

	conf, err := c.service.GetStockConfiguration(ctx, in.GetStockId())
	if err != nil {
		return nil, err
	}

	var (
		relayable api.OptionStatus
		transmittable api.OptionStatus
		storeable api.OptionStatus
	)

	if conf.Relayable == true {
		relayable = api.OptionStatus{OptionStatus: 1}
	} else {
		relayable = api.OptionStatus{OptionStatus: 0}
	}

	if conf.Transmittable == true {
		transmittable = api.OptionStatus{OptionStatus: 1}
	} else {
		transmittable = api.OptionStatus{OptionStatus: 0}
	}

	if conf.Storeable == true {
		storeable = api.OptionStatus{OptionStatus: 1}
	} else {
		storeable = api.OptionStatus{OptionStatus: 0}
	}

	return &api.StockConfig{
		StockId: conf.StockId,
		Relayable: &relayable,
		Transmittable: &transmittable,
		Storeable: &storeable,
	}, err
}


func (c *Adapter) GetStockConfigAll(ctx context.Context, in *api.Null) (*api.StockConfigList, error) {

	confList, err := c.service.GetAllStockConfiguration(ctx)
	if err != nil {
		return nil, err
	}

	resultList := make([]*api.StockConfig, len(confList))

	for i, conf := range confList {

		var (
			result api.StockConfig

			relayable api.OptionStatus
			transmittable api.OptionStatus
			storeable api.OptionStatus
		)

		if conf.Relayable == true {
			relayable = api.OptionStatus{OptionStatus: 1}
		} else {
			relayable = api.OptionStatus{OptionStatus: 0}
		}
	
		if conf.Transmittable == true {
			transmittable = api.OptionStatus{OptionStatus: 1}
		} else {
			transmittable = api.OptionStatus{OptionStatus: 0}
		}
	
		if conf.Storeable == true {
			storeable = api.OptionStatus{OptionStatus: 1}
		} else {
			storeable = api.OptionStatus{OptionStatus: 0}
		}

		result = api.StockConfig{
			StockId: conf.StockId,
			Relayable: &relayable,
			Transmittable: &transmittable,
			Storeable: &storeable,
		}

		resultList[i] = &result		
	}

	return &api.StockConfigList{StockConfig: resultList}, nil
}