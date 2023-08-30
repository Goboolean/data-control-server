package mongo_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/Goboolean/fetch-server/internal/infrastructure/mongo"
)

func Test_InsertStockBatch(t *testing.T) {

	var (
		stockId    = "stock.goboolean.test"
		stockBatch = []*mongo.StockAggregate{{}, {}, {}}
	)

	tx, err := instance.NewTx(context.Background())
	if err != nil {
		t.Errorf("failed to start transaction: %v", err)
	}

	if err := queries.InsertStockBatch(tx, stockId, stockBatch); err != nil {
		t.Errorf("failed to insert: %v", err)
	}

	if err := tx.Commit(); err != nil {
		t.Errorf("failed to commit transaction: %v", err)
	}
}

func isEqual(send, received []*mongo.StockAggregate) bool {
	if len(send) != len(received) {
		return false
	}
	for idx := range send {
		if send[idx] != received[idx] {
			return false
		}
	}
	return true
}

func Test_FetchAllStockBatch(t *testing.T) {

	var stockId = "stock.goboolean.test"

	tx, err := instance.NewTx(context.Background())
	if err != nil {
		t.Errorf("failed to start transaction: %v", err)
		return
	}

	result, err := queries.FetchAllStockBatch(tx, stockId)
	if err != nil {
		t.Errorf("FetchAllStockBatch() failed: %v", err)
		return
	}

	if err := tx.Commit(); err != nil {
		t.Errorf("failed to commit transaction: %v", err)
		return
	}

	if len(result) == 0 {
		t.Errorf("FetchAllStockBatch() failed: result is empty")
		return
	}
}

func Test_FetchAllStockBatchMassive(t *testing.T) {

	var (
		stockId    = "stock.goboolean.test"
		stockBatch = []*mongo.StockAggregate{{}, {}, {}}
	)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	tx, err := instance.NewTx(ctx)
	if err != nil {
		t.Errorf("failed to start transaction: %v", err)
	}

	stockChan := make(chan *mongo.StockAggregate, 100)

	if err := queries.FetchAllStockBatchMassive(tx, stockId, stockChan); err != nil {
		t.Errorf("FetchAllStockBatchMassive() failed: %v", err)
	}

	received := make([]*mongo.StockAggregate, 0)

loop:
	for {
		select {
		case <-ctx.Done():
			t.Errorf("FetchAllStockBatchMassive() failed with timeout")
			break loop
		case stock := <-stockChan:
			if reflect.DeepEqual(stock, &mongo.StockAggregate{}) {
				break loop
			}
			received = append(received, stock)
		}
	}

	if isEqual(stockBatch, received) {
		t.Errorf("FetchAllStockBatchMassive() failed: send and received are not equal")
	}
}
