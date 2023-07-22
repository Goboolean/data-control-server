package redis

import (
	"github.com/Goboolean/shared/pkg/resolver"
	"google.golang.org/protobuf/proto"
)



func (q *Queries) InsertStockData(tx resolver.Transactioner, stock string, stockItem *StockAggregate) error {

	data, err := proto.Marshal(stockItem)

	if err != nil {
		return err
	}

	return q.rds.client.RPush(tx.Context(), stock, &data).Err()
}

func (q *Queries) InsertStockDataBatch(tx resolver.Transactioner, stock string, stockBatch []*StockAggregate) error {

	dataBatch := make([]interface{}, len(stockBatch))

	for idx := range dataBatch {
		data, err := proto.Marshal(stockBatch[idx])

		if err != nil {
			return err
		}

		dataBatch[idx] = data
	}

	return q.rds.client.RPush(tx.Context(), stock, dataBatch...).Err()
}
