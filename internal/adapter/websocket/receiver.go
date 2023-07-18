package websocket

import (
	"github.com/Goboolean/fetch-server/internal/domain/entity"
	"github.com/Goboolean/fetch-server/internal/infrastructure/prometheus"
	"github.com/Goboolean/fetch-server/internal/infrastructure/ws"
)

func (s *StockWsAdapter) toDomainEntity(agg *ws.StockAggregate) (*entity.StockAggregateForm, error) {
	stockId, ok := s.symbolToId[agg.Symbol]
	if !ok {
		return nil, ErrStockNotFound
	}

	return &entity.StockAggregateForm{
		StockID: stockId,
		StockAggregate: entity.StockAggregate{
			Average: agg.Average,
			Min: agg.Min,
			Max: agg.Max,
			Start: agg.Start,
			End: agg.End,
			StartTime: agg.StartTime,
			EndTime: agg.EndTime,
		},
	}, nil
}


func (s *StockWsAdapter) OnReceiveStockAggs(agg *ws.StockAggregate) error {
	prometheus.DomesticStockCounter.Inc()

	data, err := s.toDomainEntity(agg)
	if err != nil {
		return err
	}

	return s.port.PlaceStockFormBatch([]*entity.StockAggregateForm{data})
}

func (s *StockWsAdapter) OnReceiveStockAggsBatch(aggs []*ws.StockAggregate) error {
	prometheus.DomesticStockCounter.Add(float64(len(aggs)))

	batch := make([]*entity.StockAggregateForm, len(aggs))

	for _, agg := range aggs {
		data, err := s.toDomainEntity(agg)
		if err != nil {
			return err
		}

		batch = append(batch, data)
	}

	return s.port.PlaceStockFormBatch(batch)
}