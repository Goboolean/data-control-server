package config

import (
	"sync"

	api "github.com/Goboolean/fetch-server/api/grpc"
	"github.com/Goboolean/fetch-server/internal/domain/port/in"
)



type StockConfiguratorAdapter struct {
	service in.ConfiguratorPort
	api.UnimplementedStockConfiguratorServer
}

var (
	instance *StockConfiguratorAdapter
	once     sync.Once
)

func New(s in.ConfiguratorPort) *StockConfiguratorAdapter {

	once.Do(func() {
		instance = &StockConfiguratorAdapter{service: s}
	})

	return instance
}
