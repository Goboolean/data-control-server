package store_test

import (
	"context"
	"testing"

	"github.com/Goboolean/fetch-server/internal/domain/service/store"
)




func Test_Store(t *testing.T) {

	var stockId string = "test"

	ctx := context.Background()

	instance := store.New(ctx)

	t.Run("StoreStock", func(t *testing.T) {
		if exists := instance.StockExists(stockId); exists {
			t.Errorf("stock exists before storing")
			return
		}
	
		if err := instance.StoreStock(stockId); err != nil {
			t.Errorf("failed to store stock: %v", err)
			return
		}
	
		if exists := instance.StockExists(stockId); !exists {
			t.Errorf("stock not exists after storing")
			return
		}
	})

	t.Run("UnstoreStock", func(t *testing.T) {
		if err := instance.UnstoreStock(stockId); err != nil {
			t.Errorf("failed to unstore stock: %v", err)
			return
		}
	
		if exists := instance.StockExists(stockId); exists {
			t.Errorf("stock exists before storing")
			return
		}	
	})
}