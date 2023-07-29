package grpc

import (
	"context"

	api "github.com/Goboolean/fetch-server/api/grpc"
)

type MockAdapter struct {
	api.UnimplementedStockConfiguratorServer
}

func NewMockAdapter() api.StockConfiguratorServer {
	return &MockAdapter{}
}

func (a *MockAdapter) UpdateStockConfigOne(context.Context, *api.StockConfig) (*api.ReturnMessage, error) {
	return &api.ReturnMessage{}, nil
}

func (a *MockAdapter) UpdateStockConfigMany(context.Context, *api.StockConfigList) (*api.ReturnMessageList, error) {
	return &api.ReturnMessageList{}, nil
}

func (a *MockAdapter) GetStockConfigOne(context.Context, *api.StockId) (*api.StockConfig, error) {
	return &api.StockConfig{}, nil
}

func (a *MockAdapter) GetStockConfigAll(context.Context, *api.Null) (*api.StockConfigList, error) {
	return &api.StockConfigList{}, nil
}
