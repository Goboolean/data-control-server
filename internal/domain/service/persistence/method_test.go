package persistence_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Goboolean/fetch-server/cmd/inject"
	cache_adapter "github.com/Goboolean/fetch-server/internal/adapter/cache"
	persistence_adapter "github.com/Goboolean/fetch-server/internal/adapter/persistence"
	"github.com/Goboolean/fetch-server/internal/adapter/transaction"
	"github.com/Goboolean/fetch-server/internal/adapter/websocket"
	"github.com/Goboolean/fetch-server/internal/domain/port/out"
	"github.com/Goboolean/fetch-server/internal/domain/service/persistence"
	"github.com/Goboolean/fetch-server/internal/infrastructure/ws/mock"
)


var (
	instance *persistence.Manager
	db out.StockPersistencePort
	cache out.StockPersistenceCachePort
)



func SetUp() {
	var err error

	ws := websocket.NewMockAdapter().(*websocket.MockAdapter)
	f := mock.New(time.Millisecond * 10, ws)
	if err := ws.RegisterFetcher(f); err != nil {
		panic(err)
	}
	relayer, err := inject.InitMockRelayer(ws)
	if err != nil {
		panic(err)
	}
	ws.RegisterReceiver(relayer)

	tx      := transaction.NewMock()
	db       = persistence_adapter.NewMockAdapter()
	cache    = cache_adapter.NewMockAdapter()
	instance, err = persistence.New(tx, db, cache, relayer, persistence.Option{BatchSize: 1})
	if err != nil {
		panic(err)
	}
	if err := relayer.FetchStock(context.Background(), "stock.facebook.usa"); err != nil {
		panic(err)
	}
}


func TearDown() {
	instance.Close()
}


func TestMain(m *testing.M) {
	SetUp()
	code := m.Run()
	TearDown()
	os.Exit(code)
}



func Test_Persistence(t *testing.T) {

	stockId := "stock.facebook.usa"

	t.Run("SubscribeRelayer", func(t *testing.T) {

		countBefore := db.(*persistence_adapter.MockAdapter).GetStoredStockCount(stockId)

		if err := instance.SubscribeRelayer(context.Background(), stockId); err != nil {
			t.Errorf("SubscribeRelayer() = %v", err)
			return
		}

		time.Sleep(100 * time.Millisecond)

		countAfter := db.(*persistence_adapter.MockAdapter).GetStoredStockCount(stockId)
		
		if diff := countAfter - countBefore; diff == 0 {
			t.Errorf("SubscribeRelayer() = %v, expected = many", diff)
			return
		}
	})


	t.Run("UnsubscribeRelayer", func(t *testing.T) {
		if err := instance.UnsubscribeRelayer(stockId); err != nil {
			t.Errorf("UnsubscribeRelayer() = %v", err)
			return
		}

		countBefore := db.(*persistence_adapter.MockAdapter).GetStoredStockCount(stockId)

		time.Sleep(100 * time.Millisecond)

		countAfter := db.(*persistence_adapter.MockAdapter).GetStoredStockCount(stockId)
	
		if diff := countAfter - countBefore; diff != 0 {
			t.Errorf("UnsubscribeRelayer() = %v, expected = %v", diff, 0)
			return
		}
	})

	t.Run("IsStockStorable (case:false)", func(t *testing.T) {
		if flag := instance.IsStockStoreable(stockId); flag {
			t.Errorf("IsStockStoreable() = %v, expected = false", flag)
			return
		}

		if err := instance.SubscribeRelayer(context.Background(), stockId); err != nil {
			t.Errorf("SubscribeRelayer() = %v", err)
			return
		}
	})

	t.Run("IsStockStorable (case:true)", func(t *testing.T) {
		if flag := instance.IsStockStoreable(stockId); !flag {
			t.Errorf("IsStockStoreable() = %v, expected = true", flag)
			return
		}

		if err := instance.UnsubscribeRelayer(stockId); err != nil {
			t.Errorf("UnsubscribeRelayer() = %v", err)
			return
		}
	})

}



func Test_WithCache(t *testing.T) {

	stockId := "stock.facebook.usa"


	var syncCount = 15
	t.Run(fmt.Sprintf("SyncCount=%d", syncCount), func(t *testing.T) {

		db.(*persistence_adapter.MockAdapter).Clear()
		cache.(*cache_adapter.MockAdapter).Clear()

		instance.SetSyncCount(syncCount)
		defer instance.SetSyncCount(0)

		if err := instance.SubscribeRelayer(context.Background(), stockId); err != nil {
			t.Errorf("SubscribeRelayer() = %v", err)
			return
		}

		defer func() {
			if err := instance.UnsubscribeRelayer(stockId); err != nil {
				t.Errorf("UnsubscribeRelayer() = %v", err)
				return
			}
		}()

		storedBefore := db.(*persistence_adapter.MockAdapter).GetStoredStockCount(stockId)
		cachedBefore := cache.(*cache_adapter.MockAdapter).GetStoredStockCount(stockId)

		time.Sleep(100 * time.Millisecond)

		storedAfter := db.(*persistence_adapter.MockAdapter).GetStoredStockCount(stockId)
		if stored := storedAfter - storedBefore; stored != 0 {
			t.Errorf("SubscribeRelayer() = %v, expected = 0", stored)
			return
		}

		cachedAfter := cache.(*cache_adapter.MockAdapter).GetStoredStockCount(stockId)
		if cached := cachedAfter - cachedBefore; cached == 0 {
			t.Errorf("SubscribeRelayer() = %v, expected = many", cached)
			return
		}

		time.Sleep(200 * time.Millisecond)

		storedFinal := db.(*persistence_adapter.MockAdapter).GetStoredStockCount(stockId)
		if stored := storedFinal - storedAfter; stored == 0 {
			t.Errorf("SubscribeRelayer() = %v, expected = many", stored)
			return
		}
	})


	var syncDuration = 150 * time.Millisecond
	t.Run(fmt.Sprintf("SyncDuration=%dms", syncDuration / time.Millisecond), func(t *testing.T) {

		db.(*persistence_adapter.MockAdapter).Clear()
		cache.(*cache_adapter.MockAdapter).Clear()

		instance.SetSyncDuration(syncDuration)
		defer instance.SetSyncDuration(0)

		if err := instance.SubscribeRelayer(context.Background(), stockId); err != nil {
			t.Errorf("SubscribeRelayer() = %v", err)
			return
		}

		defer func() {
			if err := instance.UnsubscribeRelayer(stockId); err != nil {
				t.Errorf("UnsubscribeRelayer() = %v", err)
				return
			}
		}()

		storedBefore := db.(*persistence_adapter.MockAdapter).GetStoredStockCount(stockId)
		cachedBefore := cache.(*cache_adapter.MockAdapter).GetStoredStockCount(stockId)

		time.Sleep(100 * time.Millisecond)

		storedAfter := db.(*persistence_adapter.MockAdapter).GetStoredStockCount(stockId)
		if stored := storedAfter - storedBefore; stored != 0 {
			t.Errorf("SubscribeRelayer() = %v, expected = 0", stored)
			return
		}

		cachedAfter := cache.(*cache_adapter.MockAdapter).GetStoredStockCount(stockId)
		if cached := cachedAfter - cachedBefore; cached == 0 {
			t.Errorf("SubscribeRelayer() = %v, expected = many", cached)
			return
		}

		time.Sleep(200 * time.Millisecond)
		storedFinal := db.(*persistence_adapter.MockAdapter).GetStoredStockCount(stockId)

		if stored := storedFinal - storedAfter; stored == 0 {
			t.Errorf("SubscribeRelayer() = %v, expected = many", stored)
			return
		}
	})

}

