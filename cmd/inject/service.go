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
	"github.com/Goboolean/fetch-server/internal/infrastructure/broker"
	"github.com/Goboolean/fetch-server/internal/infrastructure/rdbms"
	"github.com/Goboolean/shared/pkg/mongo"

	"github.com/Goboolean/fetch-server/internal/domain/port"
	"github.com/Goboolean/fetch-server/internal/domain/port/in"
	"github.com/Goboolean/fetch-server/internal/domain/port/out"
	"github.com/Goboolean/fetch-server/internal/domain/service/config"
	"github.com/Goboolean/fetch-server/internal/domain/service/persistence"
	"github.com/Goboolean/fetch-server/internal/domain/service/relay"
	"github.com/Goboolean/fetch-server/internal/domain/service/transmission"
	"github.com/Goboolean/fetch-server/internal/infrastructure/grpc"
	"github.com/Goboolean/fetch-server/internal/infrastructure/redis"
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

func InitMockRelayer(out.RelayerPort) (*relay.Manager, error) {
	wire.Build(MockAdapterSet, relay.New)
	return &relay.Manager{}, nil
}

func InitMockPersister(*relay.Manager, persistence.Option) (*persistence.Manager, error) {
	wire.Build(MockAdapterSet, persistence.New)
	return &persistence.Manager{}, nil
}

func InitMockTransmitter(*relay.Manager, transmission.Option) (*transmission.Manager, error) {
	wire.Build(MockAdapterSet, transmission.New)
	return &transmission.Manager{}, nil
}

func InitMockConfigurator(*relay.Manager, *persistence.Manager, *transmission.Manager) (*config.Manager, error) {
	wire.Build(MockAdapterSet, config.New)
	return &config.Manager{}, nil
}

func InitTransactor(mongo *mongo.DB, psql *rdbms.PSQL) port.TX {
	wire.Build(transaction.New)
	return &transaction.Tx{}
}

func InitRelayer(port.TX, *mongo.Queries, *rdbms.Queries, out.RelayerPort) (*relay.Manager, error) {
	wire.Build(AdapterSet, relay.New)
	return &relay.Manager{}, nil
}

func InitTransmitter(port.TX, transmission.Option, *broker.Configurator, *broker.Publisher, *relay.Manager) (*transmission.Manager, error) {
	wire.Build(AdapterSet, transmission.New)
	return &transmission.Manager{}, nil
}

func InitPersister(port.TX, persistence.Option, *redis.Queries, *rdbms.Queries, *mongo.Queries, *relay.Manager) (*persistence.Manager, error) {
	wire.Build(AdapterSet, persistence.New)
	return &persistence.Manager{}, nil
}

func InitConfigurator(port.TX, *rdbms.Queries, *persistence.Manager, *transmission.Manager, *relay.Manager) (*config.Manager, error) {
	wire.Build(AdapterSet, config.New)
	return &config.Manager{}, nil
}

func InitGrpcWithAdapter(in.ConfiguratorPort) (*grpc.Host, error) {
	wire.Build(AdapterSet, GrpcSet)
	return &grpc.Host{}, nil
}

func InitWs() *websocket.Adapter {
	wire.Build(AdapterSet)
	return &websocket.Adapter{}
}
