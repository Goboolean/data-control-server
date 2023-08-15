//go:build wireinject
// +build wireinject

package inject

import (
	broker_adapter "github.com/Goboolean/fetch-server/internal/adapter/broker"
	"github.com/Goboolean/fetch-server/internal/adapter/cache"
	grpc_adapter "github.com/Goboolean/fetch-server/internal/adapter/grpc"
	"github.com/Goboolean/fetch-server/internal/adapter/meta"
	persistence_adapter "github.com/Goboolean/fetch-server/internal/adapter/persistence"
	"github.com/Goboolean/fetch-server/internal/adapter/transaction"
	"github.com/Goboolean/fetch-server/internal/adapter/websocket"
	"github.com/Goboolean/shared/pkg/broker"
	"github.com/Goboolean/shared/pkg/mongo"
	"github.com/Goboolean/shared/pkg/rdbms"

	"github.com/Goboolean/fetch-server/internal/domain/port"
	"github.com/Goboolean/fetch-server/internal/domain/port/in"
	"github.com/Goboolean/fetch-server/internal/domain/port/out"
	"github.com/Goboolean/fetch-server/internal/domain/service/config"
	"github.com/Goboolean/fetch-server/internal/domain/service/persistence"
	"github.com/Goboolean/fetch-server/internal/domain/service/relayer"
	"github.com/Goboolean/fetch-server/internal/domain/service/transmission"
	"github.com/Goboolean/fetch-server/internal/infrastructure/cache/redis"
	"github.com/Goboolean/fetch-server/internal/infrastructure/grpc"
	"github.com/Goboolean/fetch-server/internal/infrastructure/ws"
	"github.com/google/wire"
)



func provideFetcher() []ws.Fetcher {
	return []ws.Fetcher{}
}

var AdapterArgumentSet = wire.NewSet(
	provideFetcher,
)


var MockAdapterSet = wire.NewSet(
	AdapterArgumentSet,
	grpc_adapter.NewAdapter,
	broker_adapter.NewMockAdapter,
	cache.NewMockAdapter,
	meta.NewMockAdapter,
	persistence_adapter.NewMockAdapter,
	transaction.NewMock,
	websocket.NewAdapter,
)


var AdapterSet = wire.NewSet(
	AdapterArgumentSet,
	broker_adapter.NewAdapter,
	grpc_adapter.NewAdapter,
	cache.NewAdapter,
	meta.NewAdapter,
	persistence_adapter.NewAdapter,
	websocket.NewAdapter,
)




func InitMockRelayer(out.RelayerPort) (*relayer.RelayerManager, error) {
	wire.Build(MockAdapterSet, relayer.New)
	return &relayer.RelayerManager{}, nil
}

func InitMockPersistenceManager(*relayer.RelayerManager, persistence.Option) (*persistence.PersistenceManager, error) {
	wire.Build(MockAdapterSet, persistence.New)
	return &persistence.PersistenceManager{}, nil
}

func InitMockTransmissionManager(*relayer.RelayerManager, transmission.Option) (*transmission.Transmitter, error) {
	wire.Build(MockAdapterSet, transmission.New)
	return &transmission.Transmitter{}, nil
}

func InitMockConfigurationManager(*relayer.RelayerManager, *persistence.PersistenceManager, *transmission.Transmitter) (*config.ConfigurationManager, error) {
	wire.Build(MockAdapterSet, config.New)
	return &config.ConfigurationManager{}, nil
}



func InitTransactor(mongo *mongo.DB, psql *rdbms.PSQL) port.TX {
	wire.Build(transaction.New)
	return &transaction.Tx{}
}

func InitRelayer(port.TX, *mongo.Queries, *rdbms.Queries, out.RelayerPort) (*relayer.RelayerManager, error) {
	wire.Build(AdapterSet, relayer.New)
	return &relayer.RelayerManager{}, nil
}

func InitTransmission(port.TX, transmission.Option, *broker.Configurator, *broker.Publisher, *relayer.RelayerManager) (*transmission.Transmitter, error) {
	wire.Build(AdapterSet, transmission.New)
	return &transmission.Transmitter{}, nil
}

func InitPersistenceManager(port.TX, persistence.Option, *redis.Queries, *rdbms.Queries, *mongo.Queries, *relayer.RelayerManager) (*persistence.PersistenceManager, error) {
	wire.Build(AdapterSet, persistence.New)
	return &persistence.PersistenceManager{}, nil
}

func InitConfigurationManager(port.TX, *rdbms.Queries, *persistence.PersistenceManager, *transmission.Transmitter, *relayer.RelayerManager) (*config.ConfigurationManager, error) {
	wire.Build(AdapterSet, config.New)
	return &config.ConfigurationManager{}, nil
}


func InitGrpcWithAdapter(in.ConfiguratorPort) (*grpc.Host, error) {
	wire.Build(AdapterSet, GrpcSet)
	return &grpc.Host{}, nil
}

func InitWs() *websocket.Adapter {
	wire.Build(AdapterSet)
	return &websocket.Adapter{}
}