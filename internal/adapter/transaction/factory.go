package transaction

import (
	"os"
	"sync"

	"github.com/Goboolean/fetch-server/internal/infrastructure/rediscache"
	"github.com/Goboolean/shared/pkg/mongo"
	"github.com/Goboolean/shared/pkg/rdbms"
	"github.com/Goboolean/shared/pkg/resolver"
)

type Factory struct {
	m *mongo.DB
	p *rdbms.PSQL
	r *rediscache.Redis
}

var (
	factory *Factory
	once    sync.Once
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

			m: mongo.NewDB(&resolver.ConfigMap{
				"HOST":     os.Getenv("MONGO_HOST"),
				"PORT":     os.Getenv("MONGO_PORT"),
				"USER":     os.Getenv("MONGO_USER"),
				"PASSWORD": os.Getenv("MONGO_PASS"),
			}),

			p: rdbms.NewDB(&resolver.ConfigMap{
				"HOST":     os.Getenv("PSQL_HOST"),
				"PORT":     os.Getenv("PSQL_PORT"),
				"USER":     os.Getenv("PSQL_USER"),
				"PASSWORD": os.Getenv("PSQL_PASS"),
			}),
		}
	})

	return factory
}
