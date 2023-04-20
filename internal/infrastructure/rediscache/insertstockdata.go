package rediscache

import (
	"github.com/Goboolean/data-control-server/internal/infrastructure/transaction"
	"go.mongodb.org/mongo-driver/bson"
)




func (q *Queries) InsertStockData(tx infra.Transactioner, stock string, stockItem *StockAggregate) error {

	data, err := bson.Marshal(&stockItem)

	if err != nil {
		return err
	}

	return q.rds.RPush(tx.Context(), stock, &data).Err()
}


func (q *Queries) InsertStockDataBatch(tx infra.Transactioner, stock string, stockBatch []StockAggregate) error {

	dataBatch := make([]interface{}, len(stockBatch))

	for idx := range dataBatch {
		data, err := bson.Marshal(&dataBatch[idx])

		if err != nil {
			return err
		}

		dataBatch[idx] = data
	}

	return q.rds.RPush(tx.Context(), stock, dataBatch...).Err()
}