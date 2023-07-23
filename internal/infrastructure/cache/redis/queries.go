package redis

import (
	"github.com/Goboolean/shared/pkg/resolver"
	"github.com/go-redis/redis/v8"
	"google.golang.org/protobuf/proto"
)



type Queries struct {
	db *Redis
}

func New(db *Redis) *Queries {
	return &Queries{db: db}
}



func (q *Queries) InsertStockData(tx resolver.Transactioner, stock string, stockItem *StockAggregate) error {

	data, err := proto.Marshal(stockItem)

	if err != nil {
		return err
	}

	result := q.db.client.RPush(tx.Context(), stock, data)
	return result.Err()
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

	return q.db.client.RPush(tx.Context(), stock, dataBatch...).Err()
}


func (q *Queries) GetAndEmptyCache(tx resolver.Transactioner, stock string) ([]*StockAggregate, error) {

	pipe := tx.Transaction().(*redis.Pipeline)

	getListCmd := pipe.LRange(tx.Context(), stock, 0, -1)
	getLenCmd := pipe.LLen(tx.Context(), stock)
	pipe.Del(tx.Context(), stock)

	_, err := pipe.Exec(tx.Context())

	if err != nil {
		return nil, err
	}

	length, _ := getLenCmd.Result()
	data, _ := getListCmd.Result()

	stockBatch := make([]*StockAggregate, length)

	for idx := range data {
		var stockItem StockAggregate

		if err := proto.Unmarshal([]byte(data[idx]), &stockItem); err != nil {
			return nil, err
		}
		stockBatch[idx] = &stockItem
	}

	return stockBatch, nil
}
