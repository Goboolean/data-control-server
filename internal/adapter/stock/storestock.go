package stock

import (
	adapter "github.com/Goboolean/stock-fetch-server/internal/adapter/transaction"
	"github.com/Goboolean/stock-fetch-server/internal/domain/port"
	"github.com/Goboolean/stock-fetch-server/internal/domain/value"
	"github.com/Goboolean/stock-fetch-server/internal/infrastructure/mongodb"
)

func (a *StockAdapter) StoreStock(tx port.Transactioner, stockId string, stockData []value.StockAggregate) error {
	q := mongodb.New()

	dataBatch := make([]mongodb.StockAggregate, 0)

	for idx := range stockData {
		dataBatch = append(dataBatch, mongodb.StockAggregate{
			EventType: stockData[idx].EventType,
			Avg:       stockData[idx].Average,
			Min:       stockData[idx].Min,
			Max:       stockData[idx].Max,
			Start:     stockData[idx].Start,
			End:       stockData[idx].End,
			StartTime: stockData[idx].StartTime,
			EndTime:   stockData[idx].EndTime,
		})
	}

	if err := q.InsertStockBatch(tx.(*adapter.Transaction).Mongo, stockId, dataBatch); err != nil {
		return err
	}

	return nil
}
