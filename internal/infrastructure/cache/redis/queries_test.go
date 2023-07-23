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

	length, err := queries.GetStockBatchStoredLength(ctx, stockId)
	if err != nil {
		t.Errorf("failed to get stock batch stored length: %v", err)
		return
	}

	if err := queries.InsertStockData(ctx, stockId, stockData); err != nil {
		t.Errorf("failed to insert stock data: %v", err)
		return
	}

	lengthUpdated, err := queries.GetStockBatchStoredLength(ctx, stockId)
	if err != nil {
		t.Errorf("failed to get stock batch stored length: %v", err)
		return
	}

	if lengthUpdated != length + 1 {
		t.Errorf("failed to insert stock data: length = %d, lengthUpdated = %d", length, lengthUpdated)
		return
	}

	t.Logf("length = %d, lengthUpdated = %d", length, lengthUpdated)
	
}


func Test_GetAndEmptyCache(t *testing.T) {
	
	ctx := context.Background()

	if err := queries.InsertStockDataBatch(ctx, stockId, []*redis.StockAggregate{stockData}); err != nil {
		t.Errorf("failed to insert stock data: %v", err)
		return
	}

	length, err := queries.GetStockBatchStoredLength(ctx, stockId)
	if err != nil {
		t.Errorf("failed to get stock batch stored length: %v", err)
		return
	}

	batch, err := queries.GetAndEmptyCache(ctx, stockId)
	if err != nil {
		t.Errorf("failed to get and empty cache: %v", err)
		return
	}

	if len(batch) != length {
		t.Errorf("length want = %d, len(batch) received = %d", length, len(batch))
		return
	}
	

	lengthUpdated, err := queries.GetStockBatchStoredLength(ctx, stockId)
	if err != nil {
		t.Errorf("failed to get stock batch stored length: %v", err)
		return
	}

	if lengthUpdated != 0 {
		t.Errorf("failed to get and empty cache: lengthUpdated = %d", lengthUpdated)
		return
	}
}