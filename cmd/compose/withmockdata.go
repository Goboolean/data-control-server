package compose

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/Goboolean/fetch-server/cmd/inject"
	"github.com/Goboolean/fetch-server/internal/domain/service/persistence"
	"github.com/Goboolean/fetch-server/internal/domain/service/transmission"
	"github.com/Goboolean/fetch-server/internal/infrastructure/cache/redis"
	"github.com/Goboolean/fetch-server/internal/infrastructure/ws/mock"
	"github.com/Goboolean/shared/pkg/mongo"
	"github.com/Goboolean/shared/pkg/rdbms"
	"github.com/joho/godotenv"
)

func WithMockData() (err error) {
	if err := godotenv.Load(); err != nil {
		return err
	}
	// Rule:
	// Every Constructor must be deffered with Close() even if close function has no role.
	// Ends of Close() method must assure that every goroutine it holds are closed

	// Initialize Infrastructure
	pub, err := inject.InitKafkaPublisher()
	if err != nil {
		panic(err)
	}
	defer pub.Close()

	conf, err := inject.InitKafkaConfigurator()
	if err != nil {
		panic(err)
	}
	defer conf.Close()

	mongoDB, err := inject.InitMongo()
	if err != nil {
		panic(err)
	}
	defer mongoDB.Close()
	mongoQueries := mongo.New(mongoDB)

	psqlDB, err := inject.InitPsql()
	if err != nil {
		panic(err)
	}
	defer psqlDB.Close()
	psqlQueries := rdbms.NewQueries(psqlDB)

	redisDB, err := inject.InitRedis()
	if err != nil {
		panic(err)
	}
	defer redisDB.Close()
	redisQueries := redis.New(redisDB)

	transactor := inject.InitTransactor(mongoDB, psqlDB)
	
	ws := inject.InitWs()


	// Initialize Service
	relayer, err := inject.InitRelayer(transactor, mongoQueries, psqlQueries, nil)
	if err != nil {
		panic(err)
	}
	defer relayer.Close()

	transmitter, err := inject.InitTransmission(transactor, transmission.Option{}, conf, pub, relayer)
	if err != nil {
		panic(err)
	}
	defer transmitter.Close()

	persister, err := inject.InitPersistenceManager(transactor, persistence.Option{}, redisQueries, psqlQueries, mongoQueries, relayer)
	if err != nil {
		panic(err)
	}
	defer persister.Close()

	configurator, err := inject.InitConfigurationManager(transactor, psqlQueries, persister, transmitter, relayer)
	if err != nil {
		panic(err)
	}
	defer func(){}()


	// Initialize Infrastructure
	grpc, err := inject.InitGrpcWithAdapter(configurator)
	if err != nil {
		panic(err)
	}
	defer grpc.Close()

	fetcher := mock.New(time.Millisecond * 10, ws)
	defer fetcher.Close()

	if err := ws.RegisterFetcher(fetcher); err != nil {
		panic(err)
	}
	ws.RegisterReceiver(relayer)


	// Initialize util
	prom, err := inject.InitPrometheus()
	if err != nil {
		panic(err)
	}
	defer prom.Close()

	

	
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	defer func() {
		// Every fatel error will be catched here
		if panic := recover(); err != nil {
			err = panic.(error)
		}
	}()

	<- ctx.Done()

	return fmt.Errorf("signal: %v", err)
}
