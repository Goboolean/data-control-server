package rediscache

import (
	"github.com/Goboolean/stock-fetch-server/internal/infrastructure/transaction"
	"github.com/go-redis/redis/v8"
	"google.golang.org/protobuf/proto"
)

func (q *Queries) GetAndEmptyCache(tx infra.Transactioner, stock string) ([]StockAggregate, error) {

	pipe := tx.Transaction().(*redis.Pipeline)

	getListCmd := pipe.Do(tx.Context(), "LRANGE", stock, 0, -1)
	getLenCmd := pipe.LLen(tx.Context(), stock)
	pipe.Del(tx.Context(), stock)

	_, err := pipe.Exec(tx.Context())

	if err != nil {
		return nil, err
	}

	length, _ := getLenCmd.Result()
	stockBatch := make([]StockAggregate, length)

	for idx := range stockBatch {
		data, _ := getListCmd.Result()

		if err := proto.Unmarshal(data.([]byte), &stockBatch[idx]); err != nil {
			return nil, err
		}
	}

	return stockBatch, nil
}
