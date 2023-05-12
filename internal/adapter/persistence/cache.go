package stock

import (
	"github.com/Goboolean/stock-fetch-server/internal/adapter/transaction"
	"github.com/Goboolean/stock-fetch-server/internal/domain/port"
	"github.com/Goboolean/stock-fetch-server/internal/domain/value"
	"github.com/Goboolean/stock-fetch-server/internal/infrastructure/rediscache"

	"github.com/Goboolean/shared-packages/pkg/resolver"
)



func (a *StockAdapter) EmptyCache(tx port.Transactioner, stock string) ([]value.StockAggregate, error) {

	DTO, err := a.redis.GetAndEmptyCache(tx.(*transaction.Transaction).R, stock)

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



func (a *StockAdapter) InsertOnCache(tx port.Transactioner, stock string, data []value.StockAggregate) error {

	DTO := make([]rediscache.StockAggregate, len(data))
	for idx := range data {
		DTO[idx] = rediscache.StockAggregate{
			EventType: data[idx].EventType,
			Avg:       float32(data[idx].Average),
			Min:       float32(data[idx].Min),
			Max:       float32(data[idx].Max),
			Start:     float32(data[idx].Start),
			End:       float32(data[idx].End),
			StartTime: data[idx].StartTime,
			EndTime:   data[idx].EndTime,
		}
	}

	return a.redis.InsertStockDataBatch(tx.(*transaction.Transaction).R.Transaction().(resolver.Transactioner), stock, DTO)
}
