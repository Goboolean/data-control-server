package transaction

import (
	"os"
	"sync"

	"github.com/Goboolean/shared-packages/pkg/mongo"
	"github.com/Goboolean/shared-packages/pkg/rdbms"
	"github.com/Goboolean/shared-packages/pkg/resolver"
	"github.com/Goboolean/stock-fetch-server/internal/infrastructure/rediscache"
)



type Factory struct {
	m *mongo.DB
	p *rdbms.PSQL
	r *rediscache.Redis
}

var (
	factory *Factory
	once sync.Once
)

func NewFactory() *Factory {

	once.Do(func() {
		factory = &Factory{
			r: rediscache.NewInstance(&resolver.Config{
				Host:     os.Getenv("REDIS_HOST"),
				Port:     os.Getenv("REDIS_PORT"),
				User:     os.Getenv("REDIS_USER"),
				Password: os.Getenv("REDIS_PASS"),
			}),

			m: mongo.NewDB(&resolver.Config{
				Host:     os.Getenv("MONGO_HOST"),
				Port:     os.Getenv("MONGO_PORT"),
				User:     os.Getenv("MONGO_USER"),
				Password: os.Getenv("MONGO_PASS"),
			}),

			p: rdbms.NewDB(&resolver.Config{
				Host:     os.Getenv("PSQL_HOST"),
				Port:     os.Getenv("PSQL_PORT"),
				User:     os.Getenv("PSQL_USER"),
				Password: os.Getenv("PSQL_PASS"),
			}),
		}
	})

	return factory
}