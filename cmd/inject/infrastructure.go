//go:build wireinject
// +build wireinject

package inject

import (
	"os"
	"strconv"
	"time"

	"github.com/Goboolean/fetch-server/internal/domain/port/in"
	"github.com/Goboolean/fetch-server/internal/infrastructure/grpc"
	"github.com/Goboolean/fetch-server/internal/infrastructure/redis"
	"github.com/Goboolean/fetch-server/internal/infrastructure/ws"
	"github.com/Goboolean/fetch-server/internal/infrastructure/ws/buycycle"
	"github.com/Goboolean/fetch-server/internal/infrastructure/ws/kis"
	"github.com/Goboolean/fetch-server/internal/infrastructure/ws/mock"
	"github.com/Goboolean/fetch-server/internal/infrastructure/ws/polygon"
	"github.com/Goboolean/fetch-server/internal/util/prometheus"
	"github.com/Goboolean/shared/pkg/broker"
	"github.com/Goboolean/shared/pkg/mongo"
	"github.com/Goboolean/shared/pkg/rdbms"
	"github.com/Goboolean/shared/pkg/resolver"
	"github.com/google/wire"
)

func provideMongoArgs() *resolver.ConfigMap {
	return &resolver.ConfigMap{
		"HOST":     os.Getenv("MONGO_HOST"),
		"USER":     os.Getenv("MONGO_USER"),
		"PORT":     os.Getenv("MONGO_PORT"),
		"PASSWORD": os.Getenv("MONGO_PASS"),
		"DATABASE": os.Getenv("MONGO_DATABASE"),
	}
}

func providePsqlArgs() *resolver.ConfigMap {
	return &resolver.ConfigMap{
		"HOST":     os.Getenv("PSQL_HOST"),
		"USER":     os.Getenv("PSQL_USER"),
		"PORT":     os.Getenv("PSQL_PORT"),
		"PASSWORD": os.Getenv("PSQL_PASS"),
		"DATABASE": os.Getenv("PSQL_DATABASE"),
	}
}

func provideRedisArgs() *resolver.ConfigMap {
	database, err := strconv.Atoi(os.Getenv("REDIS_DATABASE"))
	if err != nil {
		panic(err)
	}

	return &resolver.ConfigMap{
		"HOST":     os.Getenv("REDIS_HOST"),
		"PORT":     os.Getenv("REDIS_PORT"),
		"USER":     os.Getenv("REDIS_USER"),
		"PASSWORD": os.Getenv("REDIS_PASS"),
		"DATABASE": database,
	}
}

func provideKafkaArgs() *resolver.ConfigMap {
	return &resolver.ConfigMap{
		"HOST": os.Getenv("KAFKA_HOST"),
		"PORT": os.Getenv("KAFKA_PORT"),
	}
}

func provideGrpcArgs() *resolver.ConfigMap {
	return &resolver.ConfigMap{
		"PORT": os.Getenv("SERVER_PORT"),
	}
}

func provideBuycycleArgs() *resolver.ConfigMap {
	return &resolver.ConfigMap{
		"HOST": os.Getenv("BUYCYCLE_HOST"),
		"PORT": os.Getenv("BUYCYCLE_PORT"),
	}
}

func provideKISArgs() *resolver.ConfigMap {
	return &resolver.ConfigMap{
		"KIS_APPKEY": os.Getenv("KIS_APPKEY"),
		"KIS_SECRET": os.Getenv("KIS_SECRET"),
	}
}

func providePolygonArgs() *resolver.ConfigMap {
	return &resolver.ConfigMap{
		"KEY": os.Getenv("POLYGON_API_KEY"),
	}
}

func providePrometheusArgs() *resolver.ConfigMap {
	return &resolver.ConfigMap{
		"PORT": os.Getenv("METRIC_PORT"),
	}
}

var MongoSet = wire.NewSet(
	provideMongoArgs,
	mongo.NewDB,
	mongo.New,
)

var PsqlSet = wire.NewSet(
	providePsqlArgs,
	rdbms.NewDB,
	rdbms.NewQueries,
)

var KafkaSet = wire.NewSet(
	provideKafkaArgs,
	broker.NewConfigurator,
	broker.NewPublisher,
	broker.NewSubscriber,
)

var RedisSet = wire.NewSet(
	provideRedisArgs,
	redis.NewInstance,
	redis.New,
)

var GrpcSet = wire.NewSet(
	provideGrpcArgs,
	grpc.New,
	grpc.NewClient,
)

var BuycycleSet = wire.NewSet(
	provideBuycycleArgs,
	buycycle.New,
)

var KISSet = wire.NewSet(
	provideKISArgs,
	kis.New,
)

var PolygonSet = wire.NewSet(
	providePolygonArgs,
	polygon.New,
)

var PrometheusSet = wire.NewSet(
	providePrometheusArgs,
	prometheus.New,
)

func InitMongo() (*mongo.DB, error) {
	wire.Build(MongoSet)
	return &mongo.DB{}, nil
}

func InitPsql() (*rdbms.PSQL, error) {
	wire.Build(PsqlSet)
	return &rdbms.PSQL{}, nil
}

func InitRedis() (*redis.Redis, error) {
	wire.Build(RedisSet)
	return &redis.Redis{}, nil
}

func InitKafkaConfigurator() (*broker.Configurator, error) {
	wire.Build(KafkaSet)
	return &broker.Configurator{}, nil
}

func InitKafkaPublisher() (*broker.Publisher, error) {
	wire.Build(KafkaSet)
	return &broker.Publisher{}, nil
}

func InitGrpc(in.ConfiguratorPort) (*grpc.Host, error) {
	wire.Build(GrpcSet, AdapterSet)
	return &grpc.Host{}, nil
}

func InitGrpcClient() (*grpc.Client, error) {
	wire.Build(GrpcSet)
	return &grpc.Client{}, nil
}

func InitBuycycle(ws.Receiver) (*buycycle.Subscriber, error) {
	wire.Build(BuycycleSet)
	return &buycycle.Subscriber{}, nil
}

func InitKIS(ws.Receiver) (*kis.Subscriber, error) {
	wire.Build(KISSet)
	return &kis.Subscriber{}, nil
}

func InitPolygon(ws.Receiver) (*polygon.Subscriber, error) {
	wire.Build(PolygonSet)
	return &polygon.Subscriber{}, nil
}

func InitMockWebsocket(time.Duration, ws.Receiver) *mock.MockFetcher {
	wire.Build(mock.New)
	return &mock.MockFetcher{}
}

func InitPrometheus() (*prometheus.Server, error) {
	wire.Build(PrometheusSet)
	return &prometheus.Server{}, nil
}
