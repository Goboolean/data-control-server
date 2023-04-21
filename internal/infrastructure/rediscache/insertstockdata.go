package rediscache

import (
	"github.com/Goboolean/data-control-server/internal/infrastructure/transaction"
	"google.golang.org/protobuf/proto"
)




func (q *Queries) InsertStockData(tx infra.Transactioner, stock string, stockItem *StockAggregate) error {

	data, err := proto.Marshal(stockItem)

	if err != nil {
		return err
	}

	return q.rds.RPush(tx.Context(), stock, &data).Err()
}


func (q *Queries) InsertStockDataBatch(tx infra.Transactioner, stock string, stockBatch []StockAggregate) error {

	dataBatch := make([]interface{}, len(stockBatch))

	for idx := range dataBatch {
		data, err := proto.Marshal(&stockBatch[idx])

		if err != nil {
			return err
		}

		dataBatch[idx] = data
	}

	return q.rds.RPush(tx.Context(), stock, dataBatch...).Err()
}