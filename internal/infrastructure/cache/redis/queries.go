package redis

import (
	"context"

	"google.golang.org/protobuf/proto"
)



type Queries struct {
	db *Redis
}

func New(db *Redis) *Queries {
	return &Queries{db: db}
}



func (q *Queries) InsertStockData(ctx context.Context, stockId string, stockItem *StockAggregate) error {

	data, err := proto.Marshal(stockItem)

	if err != nil {
		return err
	}

	result := q.db.client.RPush(ctx, stockId, data)
	return result.Err()
}


func (q *Queries) GetStockBatchStoredLength(ctx context.Context, stockId string) (int, error) {
	result := q.db.client.LLen(ctx, stockId)
	return int(result.Val()), result.Err()
}



func (q *Queries) InsertStockDataBatch(ctx context.Context, stock string, stockBatch []*StockAggregate) error {

	dataBatch := make([]interface{}, len(stockBatch))

	for idx := range dataBatch {
		data, err := proto.Marshal(stockBatch[idx])

		if err != nil {
			return err
		}

		dataBatch[idx] = data
	}

	return q.db.client.RPush(ctx, stock, dataBatch...).Err()
}


func (q *Queries) GetAndEmptyCache(ctx context.Context, stockId string) ([]*StockAggregate, error) {

	pipe := q.db.client.Pipeline()

	getListCmd := pipe.LRange(ctx, stockId, 0, -1)
	getLenCmd := pipe.LLen(ctx, stockId)
	pipe.Del(ctx, stockId)

	_, err := pipe.Exec(ctx)

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
