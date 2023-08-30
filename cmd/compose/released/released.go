package released

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Goboolean/fetch-server.v1/cmd/inject"
	"github.com/Goboolean/fetch-server.v1/internal/domain/service/persistence"
	"github.com/Goboolean/fetch-server.v1/internal/domain/service/transmission"
	"github.com/Goboolean/fetch-server.v1/internal/infrastructure/mongo"
	"github.com/Goboolean/fetch-server.v1/internal/infrastructure/rdbms"
	"github.com/Goboolean/fetch-server.v1/internal/infrastructure/redis"
	"github.com/joho/godotenv"
)

func Run() (err error) {
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

	transmitter, err := inject.InitTransmitter(transactor, transmission.Option{}, conf, pub, relayer)
	if err != nil {
		panic(err)
	}
	defer transmitter.Close()

	persister, err := inject.InitPersister(transactor, persistence.Option{}, redisQueries, psqlQueries, mongoQueries, relayer)
	if err != nil {
		panic(err)
	}
	defer persister.Close()

	configurator, err := inject.InitConfigurator(transactor, psqlQueries, persister, transmitter, relayer)
	if err != nil {
		panic(err)
	}
	defer func() {}()

	// Initialize Infrastructure
	grpc, err := inject.InitGrpcWithAdapter(configurator)
	if err != nil {
		panic(err)
	}
	defer grpc.Close()

	fetcher, err := inject.InitKIS(ws)
	if err != nil {
		panic(err)
	}

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

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	defer cancel()

	defer func() {
		// Every fatel error will be catched here
		if panic := recover(); err != nil {
			err = panic.(error)
		}
	}()

	<-ctx.Done()

	return fmt.Errorf("signal: %v", err)
}
