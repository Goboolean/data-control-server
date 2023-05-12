package stock

import (
	"sync"

	"github.com/Goboolean/shared-packages/pkg/mongo"
	"github.com/Goboolean/shared-packages/pkg/rdbms"
	"github.com/Goboolean/stock-fetch-server/internal/infrastructure/rediscache"
)



type StockAdapter struct {
	mongo *mongo.Queries
	redis *rediscache.Queries
	postgres *rdbms.Queries
}

var (
	instance *StockAdapter
	once sync.Once
)


func NewStockAdapter(mongo *mongo.Queries, redis *rediscache.Queries, postgres *rdbms.Queries) *StockAdapter {
	
	once.Do(func() {
		instance = &StockAdapter{
			mongo: mongo,
			redis: redis,
			postgres: postgres,
		}
	})

	return instance
}

