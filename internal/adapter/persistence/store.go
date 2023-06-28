package stock

import (
	"github.com/Goboolean/shared-packages/pkg/mongo"
	"github.com/Goboolean/stock-fetch-server/internal/adapter/transaction"
	"github.com/Goboolean/stock-fetch-server/internal/domain/port"
	"github.com/Goboolean/stock-fetch-server/internal/domain/value"
	"github.com/Goboolean/stock-fetch-server/internal/infrastructure/prometheus"
)

func (a *StockAdapter) StoreStock(tx port.Transactioner, stockId string, stockData []value.StockAggregate) error {

	prometheus.StoreCounter.Add(float64(len(stockData)))

	dataBatch := make([]mongo.StockAggregate, 0)

	for idx := range stockData {
		dataBatch = append(dataBatch, mongo.StockAggregate{
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

	if err := a.mongo.InsertStockBatch(tx.(*transaction.Transaction).M, stockId, dataBatch); err != nil {
		return err
	}

	return nil
}
