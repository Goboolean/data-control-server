package main

import (
	"log"

	"github.com/Goboolean/fetch-server/cmd/inject"
	"github.com/Goboolean/fetch-server/internal/domain/service/persistence"
	"github.com/Goboolean/fetch-server/internal/domain/service/transmission"
	"github.com/Goboolean/fetch-server/internal/infrastructure/cache/redis"
	"github.com/Goboolean/shared/pkg/mongo"
	"github.com/Goboolean/shared/pkg/rdbms"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	/*
	 write constructor code here
	*/

	// Initialize Infrastructure
	pub := inject.InitKafkaPublisher()
	conf := inject.InitKafkaConfigurator()
	defer conf.Close()
	defer pub.Close()

	mongoDB := inject.InitMongo()
	mongoQueries := mongo.New(mongoDB)
	defer mongoDB.Close()

	psqlDB := inject.InitPsql()
	psqlQueries := rdbms.NewQueries(psqlDB)
	defer psqlDB.Close()

	redisDB := inject.InitRedis()
	redisQueries := redis.New(redisDB)
	defer redisDB.Close()

	transactor := inject.InitTransactor(mongoDB, psqlDB)
	
	prom := inject.InitPrometheus()
	defer prom.Close()

	ws := inject.InitWs(prom)

	// Initialize Service
	relayer := inject.InitRelayer(transactor, mongoQueries, psqlQueries, nil, prom)
	defer relayer.Close()
	transmitter := inject.InitTransmission(transactor, transmission.Option{}, conf, pub, relayer, prom)
	defer transmitter.Close()
	persister := inject.InitPersistenceManager(transactor, persistence.Option{}, redisQueries, psqlQueries, mongoQueries, relayer, prom)
	defer persister.Close()
	configurator := inject.InitConfigurationManager(transactor, psqlQueries, persister, transmitter, relayer, prom)
	defer func(){}()

	// Initialize Infrastructure
	grpc := inject.InitGrpcWithAdapter(configurator)
	defer grpc.Close()

	kis := inject.InitKIS(ws)
	defer kis.Close()

	if err := ws.RegisterFetcher(kis); err != nil {
		panic(err)
	}
	ws.RegisterReceiver(relayer)
	


	// Rule:
	// Every Constructor must be deffered with Close()
	// Ends of Close() method must assure that every goroutine it holds are closed

	defer func() {
		// every fatel error will be catched here
		// call cancelFunc to cease all process
		if err := recover(); err != nil {
			log.Fatal(err)
		}
	}()

	select{}
}
