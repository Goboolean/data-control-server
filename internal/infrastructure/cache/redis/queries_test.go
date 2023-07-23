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

	tx, err := instance.NewTx(ctx)
	if err != nil {
		t.Errorf("failed to create transaction: %v", err)
	}
	defer tx.Rollback()

	if err := queries.InsertStockData(tx, stockId, stockData); err != nil {
		t.Errorf("failed to insert stock data: %v", err)
	}

	if err := tx.Commit(); err != nil {
		t.Errorf("failed to commit transaction: %v", err)
	}
}


func Test_GetAndEmptyCache(t *testing.T) {
	
	ctx := context.Background()

	tx, err := instance.NewTx(ctx)
	if err != nil {
		t.Errorf("failed to create transaction: %v", err)
	}
	defer tx.Rollback()

	if err := queries.InsertStockDataBatch(tx, stockId, []*redis.StockAggregate{stockData}); err != nil {
		t.Errorf("failed to insert stock data: %v", err)
	}

	batch, err := queries.GetAndEmptyCache(tx, stockId)
	if err != nil {
		t.Errorf("failed to get and empty cache: %v", err)
	}

	if len(batch) != 1 {
		t.Errorf("failed to get and empty cache: %v", err)
	}

	if err := tx.Commit(); err != nil {
		t.Errorf("failed to commit transaction: %v", err)
	}
}