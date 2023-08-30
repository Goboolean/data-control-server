package benchmark_test

import (
	"os"
	"testing"
	"time"

	"github.com/Goboolean/fetch-server.v1/cmd/inject"
	"github.com/Goboolean/fetch-server.v1/internal/adapter/websocket"
	"github.com/Goboolean/fetch-server.v1/internal/domain/port"
	"github.com/Goboolean/fetch-server.v1/internal/domain/service/config"
	"github.com/Goboolean/fetch-server.v1/internal/domain/service/persistence"
	"github.com/Goboolean/fetch-server.v1/internal/domain/service/relay"
	"github.com/Goboolean/fetch-server.v1/internal/domain/service/transmission"
	grpc_infra "github.com/Goboolean/fetch-server.v1/internal/infrastructure/grpc"
	"github.com/Goboolean/fetch-server.v1/internal/infrastructure/kafka"
	"github.com/Goboolean/fetch-server.v1/internal/infrastructure/mongo"
	"github.com/Goboolean/fetch-server.v1/internal/infrastructure/rdbms"
	"github.com/Goboolean/fetch-server.v1/internal/infrastructure/redis"
	mock_infra "github.com/Goboolean/fetch-server.v1/internal/infrastructure/ws/mock"
	_ "github.com/Goboolean/fetch-server.v1/internal/util/env"
)

// This package is for bemchmark test
// It tests

var (
	pub          *kafka.Publisher
	conf         *kafka.Configurator
	sub          *kafka.Subscriber
	mongoDB      *mongo.DB
	mongoQueries *mongo.Queries
	psqlDB       *rdbms.PSQL
	psqlQueries  *rdbms.Queries
	redisDB      *redis.Redis
	redisQueries *redis.Queries
	transactor   port.TX

	relayer     *relay.Manager
	transmitter *transmission.Manager
	persister   *persistence.Manager
	Manager     *config.Manager

	grpc *grpc_infra.Host
	ws   *websocket.Adapter
	mock *mock_infra.MockFetcher

	grpcClient *grpc_infra.Client
)

func SetUp() {

	var err error

	pub, err = inject.InitKafkaPublisher()
	if err != nil {
		panic(err)
	}

	conf, err = inject.InitKafkaConfigurator()
	if err != nil {
		panic(err)
	}

	mongoDB, err = inject.InitMongo()
	if err != nil {
		panic(err)
	}
	mongoQueries := mongo.New(mongoDB)

	psqlDB, err = inject.InitPsql()
	if err != nil {
		panic(err)
	}
	psqlQueries = rdbms.NewQueries(psqlDB)

	redisDB, err = inject.InitRedis()
	if err != nil {
		panic(err)
	}
	redisQueries = redis.New(redisDB)

	transactor = inject.InitTransactor(mongoDB, psqlDB)

	// Initialize Service
	relayer, err = inject.InitRelayer(transactor, mongoQueries, psqlQueries, nil)
	if err != nil {
		panic(err)
	}

	transmitter, err = inject.InitTransmitter(transactor, transmission.Option{}, conf, pub, relayer)
	if err != nil {
		panic(err)
	}

	persister, err = inject.InitPersister(transactor, persistence.Option{}, redisQueries, psqlQueries, mongoQueries, relayer)
	if err != nil {
		panic(err)
	}

	Manager, err = inject.InitConfigurator(transactor, psqlQueries, persister, transmitter, relayer)
	if err != nil {
		panic(err)
	}

	// Initialize Infrastructure
	grpc, err = inject.InitGrpcWithAdapter(Manager)
	if err != nil {
		panic(err)
	}

	ws = inject.InitWs()

	mock = inject.InitMockWebsocket(10*time.Millisecond, ws)

	if err := ws.RegisterFetcher(mock); err != nil {
		panic(err)
	}
	ws.RegisterReceiver(relayer)

	// Initialize Client
	grpcClient, err = inject.InitGrpcClient()
	if err != nil {
		panic(err)
	}
}

func TearDown() {
	// Must add defer keyword so that closing sequence follows the order of develop.go

	defer pub.Close()
	defer conf.Close()
	defer mongoDB.Close()
	defer psqlDB.Close()
	defer redisDB.Close()

	defer relayer.Close()
	defer transmitter.Close()
	defer persister.Close()
	defer func() {}()

	defer grpc.Close()
}

func TestMain(t *testing.M) {
	SetUp()
	code := t.Run()
	TearDown()
	os.Exit(code)
}

func Test_Benchmark(b *testing.B) {

	// i

}
