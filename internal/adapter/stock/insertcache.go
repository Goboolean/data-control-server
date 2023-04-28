package stock

import (
	adapter "github.com/Goboolean/stock-fetch-server/internal/adapter/transaction"
	"github.com/Goboolean/stock-fetch-server/internal/domain/port"
	"github.com/Goboolean/stock-fetch-server/internal/domain/value"
	"github.com/Goboolean/stock-fetch-server/internal/infrastructure/rediscache"
	infra "github.com/Goboolean/stock-fetch-server/internal/infrastructure/transaction"
)

func (a *StockAdapter) InsertOnCache(tx port.Transactioner, stock string, data []value.StockAggregate) error {
	q := rediscache.New()

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

	return q.InsertStockDataBatch(tx.(*adapter.Transaction).Redis.Transaction().(infra.Transactioner), stock, DTO)
}
