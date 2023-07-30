//go:build wireinject
// +build wireinject

package inject

import (
	"os"
	"strconv"

	"github.com/Goboolean/fetch-server/internal/domain/port/in"
	"github.com/Goboolean/fetch-server/internal/infrastructure/cache/redis"
	"github.com/Goboolean/fetch-server/internal/infrastructure/grpc"
	"github.com/Goboolean/fetch-server/internal/infrastructure/prometheus"
	"github.com/Goboolean/fetch-server/internal/infrastructure/ws"
	"github.com/Goboolean/fetch-server/internal/infrastructure/ws/buycycle"
	"github.com/Goboolean/fetch-server/internal/infrastructure/ws/kis"
	"github.com/Goboolean/fetch-server/internal/infrastructure/ws/polygon"
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
		"HOST": os.Getenv("GRPC_HOST"),
		"PORT": os.Getenv("GRPC_PORT"),
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
		"KEY":  os.Getenv("POLYGON_API_KEY"),
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


func InitMongo() *mongo.DB {
	wire.Build(MongoSet)
	return &mongo.DB{}
}

func InitMongoQueries() *mongo.Queries {
	wire.Build(MongoSet)
	return &mongo.Queries{}
}

func InitPsql() *rdbms.PSQL{
	wire.Build(PsqlSet)
	return &rdbms.PSQL{}
}

func InitPsqlQueries() *rdbms.Queries {
	wire.Build(PsqlSet)
	return &rdbms.Queries{}
}

func InitRedis() *redis.Redis {
	wire.Build(RedisSet)
	return &redis.Redis{}
}

func InitRedisQueries() *redis.Queries {
	wire.Build(RedisSet)
	return &redis.Queries{}
}

func InitKafkaConfigurator() *broker.Configurator {
	wire.Build(KafkaSet)
	return &broker.Configurator{}
}

func InitKafkaPublisher() *broker.Publisher {
	wire.Build(KafkaSet)
	return &broker.Publisher{}
}


func InitGrpc(in.ConfiguratorPort) *grpc.Host {
	wire.Build(GrpcSet, AdapterSet)
	return &grpc.Host{}
}

func InitBuycycle(ws.Receiver) *buycycle.Subscriber {
	wire.Build(BuycycleSet)
	return &buycycle.Subscriber{}
}

func InitKIS(ws.Receiver) *kis.Subscriber {
	wire.Build(KISSet)
	return &kis.Subscriber{}
}

func InitPolygon(ws.Receiver) *polygon.Subscriber {
	wire.Build(PolygonSet)
	return &polygon.Subscriber{}
}

func InitPrometheus() *prometheus.Server {
	wire.Build(PrometheusSet)
	return &prometheus.Server{}
}