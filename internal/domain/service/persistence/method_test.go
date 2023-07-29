package persistence_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	cache_adapter "github.com/Goboolean/fetch-server/internal/adapter/cache"
	"github.com/Goboolean/fetch-server/internal/adapter/meta"
	persistence_adapter "github.com/Goboolean/fetch-server/internal/adapter/persistence"
	"github.com/Goboolean/fetch-server/internal/adapter/transaction"
	"github.com/Goboolean/fetch-server/internal/adapter/websocket"
	"github.com/Goboolean/fetch-server/internal/domain/port/out"
	"github.com/Goboolean/fetch-server/internal/domain/service/persistence"
	"github.com/Goboolean/fetch-server/internal/domain/service/relayer"
	"github.com/Goboolean/fetch-server/internal/infrastructure/ws/mock"
)


var (
	instance *persistence.PersistenceManager
	db out.StockPersistencePort
	cache out.StockPersistenceCachePort
)

func MockRelayer() *relayer.RelayerManager {

	var (
		db           = persistence_adapter.NewMockAdapter()
		tx           = transaction.NewMock()
		meta         = meta.NewMockAdapter()
		ws = websocket.NewAdapter().(*websocket.Adapter)
		f = mock.New(context.Background(), time.Millisecond * 10, ws)
	)

	if err := ws.RegisterFetcher(f); err != nil {
		panic(err)
	}

	instance := relayer.New(context.Background(), db, tx, meta, ws)
	ws.RegisterReceiver(instance)

	return instance
}

func SetUp() {

	var (		
		tx      = transaction.NewMock()
		relayer = MockRelayer()
	)

	db       = persistence_adapter.NewMockAdapter()
	cache    = cache_adapter.NewMockAdapter()
	instance = persistence.New(tx, context.Background(), db, cache, relayer, persistence.Option{BatchSize: 1})

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

	var count = 0

	t.Run("SubscribeRelayer", func(t *testing.T) {

		if err := instance.SubscribeRelayer(context.Background(), stockId); err != nil {
			t.Errorf("SubscribeRelayer() = %v", err)
			return
		}


		time.Sleep(100 * time.Millisecond)

		sended := db.(*persistence_adapter.MockAdapter).GetStoredStockCount(stockId)

		if sended == 0 {
			t.Errorf("SubscribeRelayer() = %v, expected = many", sended)
			return
		}

		count = sended
	})


	t.Run("UnsubscribeRelayer", func(t *testing.T) {
		if err := instance.UnsubscribeRelayer(stockId); err != nil {
			t.Errorf("UnsubscribeRelayer() = %v", err)
			return
		}

		time.Sleep(100 * time.Millisecond)

		sended := db.(*persistence_adapter.MockAdapter).GetStoredStockCount(stockId)
		if sended != count {
			t.Errorf("UnsubscribeRelayer() = %v, expected = %v", sended, count)
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

		time.Sleep(100 * time.Millisecond)

		stored := db.(*persistence_adapter.MockAdapter).GetStoredStockCount(stockId)
		if stored != 0 {
			t.Errorf("SubscribeRelayer() = %v, expected = 0", stored)
			return
		}

		time.Sleep(200 * time.Millisecond)

		stored = db.(*persistence_adapter.MockAdapter).GetStoredStockCount(stockId)
		cached := cache.(*cache_adapter.MockAdapter).GetStoredStockCount(stockId)

		t.Logf("stored = %v, cached = %v", stored, cached)
		if stored == 0 {
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

		time.Sleep(100 * time.Millisecond)

		stored := db.(*persistence_adapter.MockAdapter).GetStoredStockCount(stockId)
		if stored != 0 {
			t.Errorf("SubscribeRelayer() = %v, expected = 0", stored)
			return
		}

		time.Sleep(200 * time.Millisecond)

		stored = db.(*persistence_adapter.MockAdapter).GetStoredStockCount(stockId)
		if stored == 0 {
			t.Errorf("SubscribeRelayer() = %v, expected = many", stored)
			return
		}
	})

}

