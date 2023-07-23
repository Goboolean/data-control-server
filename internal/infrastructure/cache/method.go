package cache

import (
	"context"

	"github.com/Goboolean/fetch-server/internal/infrastructure/cache/redis"
	"github.com/Goboolean/shared/pkg/mongo"
	"github.com/Goboolean/shared/pkg/resolver"
)



func (c *StockPersistenceWithCache) InsertStockBatch(tx resolver.Transactioner, stock string, batch []*redis.StockAggregate) error {
	if err := c.redis.InsertStockDataBatch(tx.Context(), stock, batch); err != nil {
		return err
	}

	return nil
}



func (a *StockPersistenceWithCache) SyncCache(ctx context.Context, stockId string) error {

	redisDTO, err := a.redis.GetAndEmptyCache(ctx, stockId)
	if err != nil {
		return err
	}

	mongoDTO := make([]*mongo.StockAggregate, 0, len(redisDTO))

	for idx, dto := range redisDTO {
		mongoDTO[idx] = &mongo.StockAggregate{
			StockID: dto.GetStockId(),
			Avg: float64(dto.GetAvg()),
			Max: float64(dto.GetMax()),
			Min: float64(dto.GetMin()),
			Start: float64(dto.GetStart()),
			End: float64(dto.GetEnd()),
			StartTime: dto.GetStartTime(),
			EndTime: dto.GetEndTime(),
		}
	}

	if err := a.mongo.InsertStockBatch(nil, stockId, mongoDTO); err != nil {
		return err
	}

	return nil
}