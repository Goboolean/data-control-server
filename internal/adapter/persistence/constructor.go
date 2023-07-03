package stock

import (
	"sync"

	"github.com/Goboolean/fetch-server/internal/infrastructure/redis"
	"github.com/Goboolean/shared/pkg/mongo"
	"github.com/Goboolean/shared/pkg/rdbms"
)

type StockAdapter struct {
	mongo    *mongo.Queries
	redis    *redis.Queries
	postgres *rdbms.Queries
}

var (
	instance *StockAdapter
	once     sync.Once
)

func NewStockAdapter(mongo *mongo.Queries, redis *redis.Queries, postgres *rdbms.Queries) *StockAdapter {

	once.Do(func() {
		instance = &StockAdapter{
			mongo:    mongo,
			redis:    redis,
			postgres: postgres,
		}
	})

	return instance
}
