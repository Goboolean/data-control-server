package persistence

import (
	"context"

	"github.com/Goboolean/fetch-server.v1/internal/domain/port"
	"github.com/Goboolean/fetch-server.v1/internal/domain/port/out"
	"github.com/Goboolean/fetch-server.v1/internal/domain/vo"
	"github.com/Goboolean/fetch-server.v1/internal/infrastructure/mongo"
	"github.com/Goboolean/fetch-server.v1/internal/infrastructure/rdbms"
	"github.com/Goboolean/fetch-server.v1/internal/util/prometheus"
)

type Adapter struct {
	rdbms *rdbms.Queries
	mongo *mongo.Queries
}

func NewAdapter(rdbms *rdbms.Queries, mongo *mongo.Queries) out.StockPersistencePort {
	return &Adapter{
		rdbms: rdbms,
		mongo: mongo,
	}
}

func (a *Adapter) StoreStock(tx port.Transactioner, stockId string, agg *vo.StockAggregate) error {
	dto := &mongo.StockAggregate{
		Min:    agg.Min,
		Max:    agg.Max,
		Open:   agg.Open,
		Closed: agg.Closed,
		Time:   agg.Time,
	}

	if err := a.mongo.InsertStockBatch(tx.(*mongo.Transaction), stockId, []*mongo.StockAggregate{dto}); err != nil {
		return err
	}

	prometheus.FetchCounter.Inc()
	return nil
}

func (a *Adapter) StoreStockBatch(tx port.Transactioner, stockId string, aggs []*vo.StockAggregate) error {
	dtos := make([]*mongo.StockAggregate, 0, len(aggs))

	for _, agg := range aggs {
		dtos = append(dtos, &mongo.StockAggregate{
			Min:    agg.Min,
			Max:    agg.Max,
			Open:   agg.Open,
			Closed: agg.Closed,
			Time:   agg.Time,
		})
	}

	if err := a.mongo.InsertStockBatch(tx.(*mongo.Transaction), stockId, dtos); err != nil {
		return err
	}

	prometheus.StoreCounter.Add(float64(len(aggs)))
	return nil
}

func (a *Adapter) CreateStoringStartedLog(ctx context.Context, stockId string) error {

	return a.rdbms.CreateAccessInfo(ctx, rdbms.CreateAccessInfoParams{
		ProductID: stockId,
		Status:    "started",
	})
}

func (a *Adapter) CreateStoringStoppedLog(ctx context.Context, stockId string) error {

	return a.rdbms.CreateAccessInfo(ctx, rdbms.CreateAccessInfoParams{
		ProductID: stockId,
		Status:    "stopped",
	})
}

func (a *Adapter) CreateStoringFailedLog(ctx context.Context, stockId string) error {

	return a.rdbms.CreateAccessInfo(ctx, rdbms.CreateAccessInfoParams{
		ProductID: stockId,
		Status:    "failed",
	})
}
