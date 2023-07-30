package cache

import (
	"sync"

	"github.com/Goboolean/fetch-server/internal/infrastructure/cache/redis"
	"github.com/Goboolean/shared/pkg/mongo"
)


// It is a infrastructure that stores stocks with cache.
// There is severe problem about transaction:
// redis does not support transaction, whereas command receives transactioner.
// For this reason, it is not recommended to use this infrastructure.
type StockPersistenceWithCache struct {
	mongo *mongo.Queries
	redis *redis.Queries
}

var (
	once sync.Once
	instance *StockPersistenceWithCache
)


func New(m *mongo.DB, r *redis.Redis) *StockPersistenceWithCache {

	once.Do(func() {
		instance = &StockPersistenceWithCache{
			mongo: mongo.New(m),
			redis: redis.New(r),
		}
	})

	return instance
}


func (c *StockPersistenceWithCache) Close() {

}