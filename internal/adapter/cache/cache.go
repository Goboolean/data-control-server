package cache

import (
	"context"

	"github.com/Goboolean/fetch-server.v1/api/model"
	"github.com/Goboolean/fetch-server.v1/internal/domain/port/out"
	"github.com/Goboolean/fetch-server.v1/internal/domain/vo"
	"github.com/Goboolean/fetch-server.v1/internal/infrastructure/redis"
)

type Adapter struct {
	redis *redis.Queries
}

func NewAdapter(r *redis.Queries) out.StockPersistenceCachePort {
	return &Adapter{
		redis: r,
	}
}

func (a *Adapter) StoreStockOnCache(ctx context.Context, stockId string, stock *vo.StockAggregate) error {

	dto := &model.StockAggregate{
		EventType: stock.EventType,
		Min:       stock.Min,
		Max:       stock.Max,
		StartTime: stock.Time,
	}

	return a.redis.InsertStockData(ctx, stockId, dto)
}

func (a *Adapter) StoreStockBatchOnCache(ctx context.Context, stockId string, stockBatch []*vo.StockAggregate) error {

	dtos := make([]*model.StockAggregate, 0, len(stockBatch))

	for _, stock := range stockBatch {
		dtos = append(dtos, &model.StockAggregate{
			EventType: stock.EventType,
			Open:      stock.Open,
			Min:       stock.Min,
			Max:       stock.Max,
			StartTime: stock.Time,
		})
	}

	return a.redis.InsertStockDataBatch(ctx, stockId, dtos)
}

func (a *Adapter) GetAndEmptyCache(ctx context.Context, stockId string) ([]*vo.StockAggregate, error) {

	dtos, err := a.redis.GetAndEmptyCache(ctx, stockId)
	if err != nil {
		return nil, err
	}

	stockBatch := make([]*vo.StockAggregate, 0, len(dtos))
	for idx, dto := range dtos {
		stockBatch[idx] = &vo.StockAggregate{
			EventType: dto.EventType,
			Min:       dto.Min,
			Max:       dto.Max,
			Open:      dto.Open,
			Closed:    dto.Closed,
			Time:      dto.StartTime,
		}
	}

	return stockBatch, nil
}
