package redis

import (
	"github.com/Goboolean/shared/pkg/resolver"
	"github.com/go-redis/redis/v8"
	"google.golang.org/protobuf/proto"
)



type Queries struct {
	rds *Redis
}

func New() *Queries {
	return &Queries{}
}



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


func (q *Queries) GetAndEmptyCache(tx resolver.Transactioner, stock string) ([]*StockAggregate, error) {

	pipe := tx.Transaction().(*redis.Pipeline)

	getListCmd := pipe.Do(tx.Context(), "LRANGE", stock, 0, -1)
	getLenCmd := pipe.LLen(tx.Context(), stock)
	pipe.Del(tx.Context(), stock)

	_, err := pipe.Exec(tx.Context())

	if err != nil {
		return nil, err
	}

	length, _ := getLenCmd.Result()
	stockBatch := make([]*StockAggregate, length)

	for idx := range stockBatch {
		data, _ := getListCmd.Result()

		if err := proto.Unmarshal(data.([]byte), stockBatch[idx]); err != nil {
			return nil, err
		}
	}

	return stockBatch, nil
}
