package stock

import (
	adapter "github.com/Goboolean/stock-fetch-server/internal/adapter/transaction"
	"github.com/Goboolean/stock-fetch-server/internal/domain/port"
	"github.com/Goboolean/stock-fetch-server/internal/domain/value"
	"github.com/Goboolean/stock-fetch-server/internal/infrastructure/rediscache"
	infra "github.com/Goboolean/stock-fetch-server/internal/infrastructure/transaction"
)

func (a *StockAdapter) EmptyCache(tx port.Transactioner, stock string) ([]value.StockAggregate, error) {
	q := rediscache.New()

	DTO, err := q.GetAndEmptyCache(tx.(*adapter.Transaction).Redis.Transaction().(infra.Transactioner), stock)

	if err != nil {
		return nil, err
	}

	stockBatch := make([]value.StockAggregate, len(DTO))

	for idx := range stockBatch {
		stockBatch[idx] = value.StockAggregate{
			EventType: DTO[idx].GetEventType(),
			Average:   float64(DTO[idx].GetAvg()),
			Min:       float64(DTO[idx].GetMin()),
			Max:       float64(DTO[idx].GetMax()),
			Start:     float64(DTO[idx].GetStart()),
			End:       float64(DTO[idx].GetStartTime()),
			StartTime: DTO[idx].GetStartTime(),
			EndTime:   DTO[idx].GetEndTime(),
		}
	}

	return stockBatch, nil
}
