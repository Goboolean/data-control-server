//go:build wireinject
// +build wireinject

package inject

import (
	"os"
	"strconv"

	"github.com/Goboolean/fetch-server/internal/infrastructure/cache/redis"
	"github.com/Goboolean/fetch-server/internal/infrastructure/grpc"
	"github.com/Goboolean/fetch-server/internal/infrastructure/prometheus"
	"github.com/Goboolean/fetch-server/internal/infrastructure/ws/buycycle"
	"github.com/Goboolean/fetch-server/internal/infrastructure/ws/kis"
	"github.com/Goboolean/fetch-server/internal/infrastructure/ws/polygon"
	"github.com/Goboolean/shared/pkg/broker"
	"github.com/Goboolean/shared/pkg/mongo"
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

func provideKafkaArgs() *resolver.ConfigMap {
	return &resolver.ConfigMap{
		"HOST": os.Getenv("KAFKA_HOST"),
		"PORT": os.Getenv("KAFKA_PORT"),
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



func InitMongo() *mongo.DB {
	wire.Build(mongo.NewDB, provideMongoArgs)
	return &mongo.DB{}
}

func InitKafkaConfigurator() *broker.Configurator {
	wire.Build(broker.NewConfigurator, provideKafkaArgs)
	return &broker.Configurator{}
}

func InitKafkaPublisher() *broker.Publisher {
	wire.Build(broker.NewPublisher, provideKafkaArgs)
	return &broker.Publisher{}
}

func InitRedis() *redis.Redis {
	wire.Build(redis.NewInstance, provideRedisArgs)
	return &redis.Redis{}
}

func InitGrpc() *grpc.Host {
	wire.Build(grpc.New, provideGrpcArgs)
	return &grpc.Host{}
}

func InitBuycycle() *buycycle.Subscriber {
	wire.Build(buycycle.New, provideBuycycleArgs)
	return &buycycle.Subscriber{}
}

func InitKIS() *kis.Subscriber {
	wire.Build(kis.New, provideKISArgs)
	return &kis.Subscriber{}
}

func InitPolygon() *polygon.Subscriber {
	wire.Build(polygon.New, providePolygonArgs)
	return &polygon.Subscriber{}
}

func InitPrometheus() *resolver.ConfigMap {
	wire.Build(prometheus.New, providePrometheusArgs)
	return &resolver.ConfigMap{}
}

var InfrastructureSet = wire.NewSet(
	InitMongo,
	InitKafkaConfigurator,
	InitKafkaPublisher,
	InitRedis,
	InitGrpc,
	InitBuycycle,
	InitKIS,
	InitPolygon,
	InitPrometheus,
)