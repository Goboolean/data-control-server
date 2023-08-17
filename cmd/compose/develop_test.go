package compose_test

import (
	"context"
	"os"
	"testing"
	"time"

	grpcapi "github.com/Goboolean/fetch-server/api/grpc"
	"github.com/Goboolean/fetch-server/cmd/inject"
	"github.com/Goboolean/fetch-server/internal/adapter/websocket"
	"github.com/Goboolean/fetch-server/internal/domain/port"
	"github.com/Goboolean/fetch-server/internal/domain/service/config"
	"github.com/Goboolean/fetch-server/internal/domain/service/persistence"
	relayer_service "github.com/Goboolean/fetch-server/internal/domain/service/relayer"
	"github.com/Goboolean/fetch-server/internal/domain/service/transmission"
	"github.com/Goboolean/fetch-server/internal/infrastructure/cache/redis"
	grpc_infra "github.com/Goboolean/fetch-server/internal/infrastructure/grpc"
	mock_infra "github.com/Goboolean/fetch-server/internal/infrastructure/ws/mock"
	_ "github.com/Goboolean/fetch-server/internal/util/env"
	"github.com/Goboolean/shared/pkg/broker"
	"github.com/Goboolean/shared/pkg/mongo"
	"github.com/Goboolean/shared/pkg/rdbms"
	"github.com/stretchr/testify/assert"
)




// This package does not test develop.go directly,
// It tests the integration of all the packages that develop.go uses.
// If all of this tests passes and develop.go is not broken, it may considered a success.


var (
	pub *broker.Publisher
	conf *broker.Configurator
	mongoDB *mongo.DB
	mongoQueries *mongo.Queries
	psqlDB *rdbms.PSQL
	psqlQueries *rdbms.Queries
	redisDB *redis.Redis
	redisQueries *redis.Queries
	transactor port.TX

	relayer *relayer_service.RelayerManager
	transmitter *transmission.Transmitter
	persister *persistence.PersistenceManager
	configurator *config.ConfigurationManager

	grpc *grpc_infra.Host
	ws *websocket.Adapter
	mock *mock_infra.MockFetcher

	grpcClient *grpc_infra.Client
)



func SetUp() {

	os.Exit(0)

	var err error

	pub, err = inject.InitKafkaPublisher()
	if err != nil {
		panic(err)
	}

	conf, err = inject.InitKafkaConfigurator()
	if err != nil {
		panic(err)
	}
	defer conf.Close()

	mongoDB, err = inject.InitMongo()
	if err != nil {
		panic(err)
	}
	defer mongoDB.Close()
	mongoQueries := mongo.New(mongoDB)

	psqlDB, err = inject.InitPsql()
	if err != nil {
		panic(err)
	}
	defer psqlDB.Close()
	psqlQueries = rdbms.NewQueries(psqlDB)

	redisDB, err = inject.InitRedis()
	if err != nil {
		panic(err)
	}
	defer redisDB.Close()
	redisQueries = redis.New(redisDB)

	transactor = inject.InitTransactor(mongoDB, psqlDB)

	// Initialize Service
	relayer, err = inject.InitRelayer(transactor, mongoQueries, psqlQueries, nil)
	if err != nil {
		panic(err)
	}
	defer relayer.Close()

	transmitter, err = inject.InitTransmission(transactor, transmission.Option{}, conf, pub, relayer)
	if err != nil {
		panic(err)
	}
	defer transmitter.Close()

	persister, err = inject.InitPersistenceManager(transactor, persistence.Option{}, redisQueries, psqlQueries, mongoQueries, relayer)
	if err != nil {
		panic(err)
	}
	defer persister.Close()

	configurator, err = inject.InitConfigurationManager(transactor, psqlQueries, persister, transmitter, relayer)
	if err != nil {
		panic(err)
	}
	defer func(){}()


	// Initialize Infrastructure
	grpc, err = inject.InitGrpcWithAdapter(configurator)
	if err != nil {
		panic(err)
	}

	ws := inject.InitWs()

	mock = inject.InitMockWebsocket(10*time.Millisecond, ws)

	if err := ws.RegisterFetcher(mock); err != nil {
		panic(err)
	}
	ws.RegisterReceiver(relayer)
}



func TearDown() {
	// Add defer keyword so that closing sequence follows the order of develop.go

	defer pub.Close()
	defer conf.Close()
	defer mongoDB.Close()
	defer psqlDB.Close()
	defer redisDB.Close()

	defer relayer.Close()
	defer transmitter.Close()
	defer persister.Close()
	defer func(){}()

	defer grpc.Close()
}



func TestMain(t *testing.M) {
	SetUp()
	code := t.Run()
	TearDown()
	os.Exit(code)
}



func Test_Integration_Configuration(t *testing.T) {
	t.Skip("not now")

	type args struct {
		stockId string
		stockConfig *grpcapi.StockConfig
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "",

		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			initial, err := grpcClient.GetStockConfigOne(context.Background(), &grpcapi.StockId{
				StockId: tt.args.stockId,
			})
			assert.NoError(t, err)

			msg, err := grpcClient.UpdateStockConfigOne(context.Background(), tt.args.stockConfig)
			assert.NoError(t, err)
			assert.Equal(t, true, msg.Status)

			reply, err := grpcClient.GetStockConfigOne(context.Background(), &grpcapi.StockId{
				StockId: tt.args.stockId,
			})

			if tt.args.stockConfig.Relayable.OptionStatus == -1 {
				assert.Equal(t, initial.Relayable.OptionStatus, reply.Relayable.OptionStatus)
			} else {
				assert.Equal(t, tt.args.stockConfig.Relayable.OptionStatus, reply.Relayable.OptionStatus)
			}

			if tt.args.stockConfig.Storeable.OptionStatus == -1 {
				assert.Equal(t, initial.Storeable.OptionStatus, reply.Storeable.OptionStatus)
			} else {
				assert.Equal(t, tt.args.stockConfig.Storeable.OptionStatus, reply.Storeable.OptionStatus)
			}

			if tt.args.stockConfig.Transmittable.OptionStatus == -1 {
				assert.Equal(t, initial.Transmittable.OptionStatus, reply.Transmittable.OptionStatus)
			} else {
				assert.Equal(t, tt.args.stockConfig.Transmittable.OptionStatus, reply.Transmittable.OptionStatus)
			}			
		})
	}
}



func Test_Integration_DataPipelining(t *testing.T) {

	var stockId string = "stock.google.usa"

	t.Run("RequestAllOptionTrue", func(t *testing.T) {

		msg, err := grpcClient.UpdateStockConfigOne(context.Background(), &grpcapi.StockConfig{
			StockId: stockId,
			Relayable: &grpcapi.OptionStatus{OptionStatus: 1},
			Storeable: &grpcapi.OptionStatus{OptionStatus: 1},
			Transmittable: &grpcapi.OptionStatus{OptionStatus: 1},
		})

		assert.NoError(t, err)
		assert.Equal(t, true, msg.Status)
	})

	t.Run("CheckStockTransmitted", func(t *testing.T) {

	})

//	t.Run("CheckStockStored", func(t *testing.T) {
//		mongoQueries.
//
//	})

	
}