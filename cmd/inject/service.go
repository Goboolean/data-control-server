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
	"github.com/Goboolean/fetch-server/internal/infrastructure/prometheus"
	"github.com/Goboolean/fetch-server/internal/infrastructure/ws"
	"github.com/google/wire"
)



var ServiceSet = wire.NewSet(
	relayer.New,
	persistence.New,
	transmission.New,
	config.New,
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

/*
func InitMockRelayer() *relayer.RelayerManager {
	wire.Build(MockAdapterSet, ServiceSet)
	return &relayer.RelayerManager{}
}

func InitMockPersistenceManager(o persistence.Option) *persistence.PersistenceManager {
	wire.Build(MockAdapterSet, ServiceSet)
	return &persistence.PersistenceManager{}
}

func InitMockTransmission(o transmission.Option) *transmission.Transmitter {
	wire.Build(MockAdapterSet, ServiceSet)
	return &transmission.Transmitter{}
}

func InitMockConfigurationManager(*relayer.RelayerManager, *persistence.PersistenceManager, *transmission.Transmitter) *config.ConfigurationManager {
	wire.Build(MockAdapterSet, config.New)
	return &config.ConfigurationManager{}
}
*/


var AdapterSet = wire.NewSet(
	AdapterArgumentSet,
	broker_adapter.NewAdapter,
	grpc_adapter.NewAdapter,
	cache.NewAdapter,
	meta.NewAdapter,
	persistence_adapter.NewAdapter,
	websocket.NewAdapter,
)

/*
func InitRelayer() *relayer.RelayerManager {
	wire.Build(AdapterSet, InfrastructureSet, ServiceSet)
	return &relayer.RelayerManager{}
}

func InitPersistenceManager(o persistence.Option) *persistence.PersistenceManager {
	wire.Build(AdapterSet, InfrastructureSet, ServiceSet)
	return &persistence.PersistenceManager{}
}

func InitTransmission(o transmission.Option) *transmission.Transmitter {
	wire.Build(AdapterSet, InfrastructureSet, ServiceSet)
	return &transmission.Transmitter{}
}

func InitConfigurationManager(*relayer.RelayerManager, *persistence.PersistenceManager, *transmission.Transmitter) *config.ConfigurationManager {
	wire.Build(AdapterSet, InfrastructureSet, config.New)
	return &config.ConfigurationManager{}
}
*/

func provideTransmissionArgs() transmission.Option {
	return transmission.Option{}
}

func providePersistenceArgs() persistence.Option {
	return persistence.Option{}
}



func InitMockRelayer(out.RelayerPort) *relayer.RelayerManager {
	wire.Build(MockAdapterSet, ServiceSet)
	return &relayer.RelayerManager{}
}

func InitMockPersistenceManager(*relayer.RelayerManager, persistence.Option) *persistence.PersistenceManager {
	wire.Build(MockAdapterSet, persistence.New)
	return &persistence.PersistenceManager{}
}

func InitMockTransmissionManager(*relayer.RelayerManager, transmission.Option) *transmission.Transmitter {
	wire.Build(MockAdapterSet, transmission.New)
	return &transmission.Transmitter{}
}

func InitMockConfigurationManager(*relayer.RelayerManager, *persistence.PersistenceManager, *transmission.Transmitter) *config.ConfigurationManager {
	wire.Build(MockAdapterSet, config.New)
	return &config.ConfigurationManager{}
}



func InitTransactor(mongo *mongo.DB, psql *rdbms.PSQL) port.TX {
	wire.Build(transaction.New)
	return &transaction.Tx{}
}

func InitRelayer(port.TX, *mongo.Queries, *rdbms.Queries, out.RelayerPort, *prometheus.Server) *relayer.RelayerManager {
	wire.Build(AdapterSet, relayer.New)
	return &relayer.RelayerManager{}
}

func InitTransmission(port.TX, transmission.Option, *broker.Configurator, *broker.Publisher, *relayer.RelayerManager, *prometheus.Server) *transmission.Transmitter {
	wire.Build(AdapterSet, transmission.New)
	return &transmission.Transmitter{}
}

func InitPersistenceManager(port.TX, persistence.Option, *redis.Queries, *rdbms.Queries, *mongo.Queries, *relayer.RelayerManager, *prometheus.Server) *persistence.PersistenceManager {
	wire.Build(AdapterSet, persistence.New)
	return &persistence.PersistenceManager{}
}

func InitConfigurationManager(port.TX, *rdbms.Queries, *persistence.PersistenceManager, *transmission.Transmitter, *relayer.RelayerManager, *prometheus.Server) *config.ConfigurationManager {
	wire.Build(AdapterSet, config.New)
	return &config.ConfigurationManager{}
}


func InitGrpcWithAdapter(in.ConfiguratorPort) *grpc.Host {
	wire.Build(AdapterSet, GrpcSet)
	return &grpc.Host{}
}

func InitWs(*prometheus.Server) *websocket.Adapter {
	wire.Build(AdapterSet)
	return &websocket.Adapter{}
}