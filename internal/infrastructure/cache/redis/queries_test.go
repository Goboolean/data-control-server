package redis_test

import (
	"context"
	"testing"

	"github.com/Goboolean/fetch-server/internal/infrastructure/cache/redis"
)



var (
	stockId = "test"
	stockData = &redis.StockAggregate{}
)


func Test_InsertStockData(t *testing.T) {

	ctx := context.Background()
	if err := queries.InsertStockData(ctx, stockId, stockData); err != nil {
		t.Errorf("failed to insert stock data: %v", err)
	}
}


func Test_GetAndEmptyCache(t *testing.T) {
	
	ctx := context.Background()

	if err := queries.InsertStockDataBatch(ctx, stockId, []*redis.StockAggregate{stockData}); err != nil {
		t.Errorf("failed to insert stock data: %v", err)
	}

	batch, err := queries.GetAndEmptyCache(ctx, stockId)
	if err != nil {
		t.Errorf("failed to get and empty cache: %v", err)
	}

	if len(batch) != 1 {
		t.Errorf("failed to get and empty cache: %v", err)
	}
}