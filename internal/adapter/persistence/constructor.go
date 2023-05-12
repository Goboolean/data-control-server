package stock

import (
	"github.com/Goboolean/shared-packages/pkg/mongo"
	"github.com/Goboolean/shared-packages/pkg/rdbms"
	"github.com/Goboolean/stock-fetch-server/internal/infrastructure/rediscache"
)




type StockAdapter struct {
	mongo *mongo.Queries
	redis *rediscache.Queries
	postgres *rdbms.Queries
}


func NewStockAdapter(mongo *mongo.Queries, redis *rediscache.Queries, postgres *rdbms.Queries) *StockAdapter {
	return &StockAdapter{
		mongo: mongo,
		redis: redis,
		postgres: postgres,
	}
}

