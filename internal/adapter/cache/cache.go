package cache

import (
	"context"

	"github.com/Goboolean/fetch-server/internal/domain/entity"
	"github.com/Goboolean/fetch-server/internal/infrastructure/cache/redis"
)



type Adapter struct {
	redis *redis.Queries
}


func NewAdapter(r *redis.Redis) *Adapter {
	return &Adapter{
		redis: redis.New(r),
	}
}




func (a *Adapter) StoreStockOnCache(ctx context.Context, stockId string, stock *entity.StockAggregate) error {

	dto := &redis.StockAggregate{
		EventType: stock.EventType,
		Avg:       stock.Average,
		Min:       stock.Min,
		Max:       stock.Max,
		Start:     stock.Start,
		End:       stock.End,
		StartTime: stock.StartTime,
		EndTime:   stock.EndTime,
	}

	return a.redis.InsertStockData(ctx, stockId, dto)
}


func (a *Adapter) StoreStockBatchOnCache(ctx context.Context, stockId string, stockBatch []*entity.StockAggregate) error {
	
	dtos := make([]*redis.StockAggregate, 0, len(stockBatch))

	for _, stock := range stockBatch {
		dtos = append(dtos, &redis.StockAggregate{
			EventType: stock.EventType,
			Avg:       stock.Average,
			Min:       stock.Min,
			Max:       stock.Max,
			Start:     stock.Start,
			End:       stock.End,
			StartTime: stock.StartTime,
			EndTime:   stock.EndTime,
		})
	}

	return a.redis.InsertStockDataBatch(ctx, stockId, dtos)
}

func (a *Adapter) GetAndEmptyCache(ctx context.Context, stockId string) ([]*entity.StockAggregate, error) {

	dtos, err := a.redis.GetAndEmptyCache(ctx, stockId)
	if err != nil {
		return nil, err
	}

	stockBatch := make([]*entity.StockAggregate, 0, len(dtos))
	for idx, dto := range dtos {
		stockBatch[idx] = &entity.StockAggregate{
			EventType: dto.EventType,
			Average:   dto.Avg,
			Min:       dto.Min,
			Max:       dto.Max,
			Start:     dto.Start,
			End:       dto.End,
			StartTime: dto.StartTime,
			EndTime:   dto.EndTime,
		}
	}

	return stockBatch, nil
}