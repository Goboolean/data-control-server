package develop_test

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
	"github.com/Goboolean/fetch-server/internal/domain/service/relay"
	"github.com/Goboolean/fetch-server/internal/domain/service/transmission"
	grpc_infra "github.com/Goboolean/fetch-server/internal/infrastructure/grpc"
	"github.com/Goboolean/fetch-server/internal/infrastructure/rdbms"
	"github.com/Goboolean/fetch-server/internal/infrastructure/redis"
	mock_infra "github.com/Goboolean/fetch-server/internal/infrastructure/ws/mock"
	_ "github.com/Goboolean/fetch-server/internal/util/env"
	"github.com/Goboolean/shared/pkg/broker"
	"github.com/Goboolean/shared/pkg/mongo"
	"github.com/stretchr/testify/assert"
)

// This package does not test develop.go directly,
// It tests the integration of all the packages that develop.go uses.
// If all of this tests passes and develop.go is not broken, it may considered a success.

var (
	pub          *broker.Publisher
	conf         *broker.Configurator
	sub          *broker.Subscriber
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

	// TODO: Check if RegisterReceiver can be first
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

func Test_Integration_Configuration(t *testing.T) {

	type args struct {
		stockConfig *grpcapi.StockConfig
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Random",
			args: args{
				stockConfig: &grpcapi.StockConfig{
					StockId:       "stock.google.usa",
					Relayable:     &grpcapi.OptionStatus{OptionStatus: 1},
					Storeable:     &grpcapi.OptionStatus{OptionStatus: 1},
					Transmittable: &grpcapi.OptionStatus{OptionStatus: 1},
				},
			},
			wantErr: false,
		},
		{
			name: "Random",
			args: args{
				stockConfig: &grpcapi.StockConfig{
					StockId:       "stock.google.usa",
					Relayable:     &grpcapi.OptionStatus{OptionStatus: 1},
					Storeable:     &grpcapi.OptionStatus{OptionStatus: 1},
					Transmittable: &grpcapi.OptionStatus{OptionStatus: 1},
				},
			},
			wantErr: false,
		},
		{
			name: "Random",
			args: args{
				stockConfig: &grpcapi.StockConfig{
					StockId:       "stock.google.usa",
					Relayable:     &grpcapi.OptionStatus{OptionStatus: -1},
					Storeable:     &grpcapi.OptionStatus{OptionStatus: 0},
					Transmittable: &grpcapi.OptionStatus{OptionStatus: 1},
				},
			},
			wantErr: false,
		},
		{
			name: "Random",
			args: args{
				stockConfig: &grpcapi.StockConfig{
					StockId:       "stock.google.usa",
					Relayable:     &grpcapi.OptionStatus{OptionStatus: 0},
					Storeable:     &grpcapi.OptionStatus{OptionStatus: 1},
					Transmittable: &grpcapi.OptionStatus{OptionStatus: 1},
				},
			},
			wantErr: false,
		},
		{
			name: "Random",
			args: args{
				stockConfig: &grpcapi.StockConfig{
					StockId:       "stock.google.usa",
					Relayable:     &grpcapi.OptionStatus{OptionStatus: 1},
					Storeable:     &grpcapi.OptionStatus{OptionStatus: -1},
					Transmittable: &grpcapi.OptionStatus{OptionStatus: 0},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			initial, err := grpcClient.GetStockConfigOne(context.Background(), &grpcapi.StockId{
				StockId: tt.args.stockConfig.StockId,
			})
			assert.NoError(t, err)

			msg, err := grpcClient.UpdateStockConfigOne(context.Background(), tt.args.stockConfig)
			assert.NoError(t, err)
			assert.Equal(t, true, msg.Status)

			reply, err := grpcClient.GetStockConfigOne(context.Background(), &grpcapi.StockId{
				StockId: tt.args.stockConfig.StockId,
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
			StockId:       stockId,
			Relayable:     &grpcapi.OptionStatus{OptionStatus: 1},
			Storeable:     &grpcapi.OptionStatus{OptionStatus: 1},
			Transmittable: &grpcapi.OptionStatus{OptionStatus: 1},
		})

		assert.NoError(t, err)
		assert.Equal(t, true, msg.Status)
	})

	t.Run("CheckStockTransmitted", func(t *testing.T) {

	})

	t.Run("CheckStockStored", func(t *testing.T) {

	})

}
