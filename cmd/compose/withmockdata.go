package compose

import (
	"fmt"
	"os"
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

	fetcher := mock.New(time.Millisecond * 10, ws)
	defer fetcher.Close()

	if err := ws.RegisterFetcher(fetcher); err != nil {
		panic(err)
	}
	ws.RegisterReceiver(relayer)
	

	
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	defer func() {
		// Every fatel error will be catched here
		if panic := recover(); err != nil {
			err = panic.(error)
		}
	}()

	sig := <- sigs
	return fmt.Errorf("signal: %v", sig)
}
