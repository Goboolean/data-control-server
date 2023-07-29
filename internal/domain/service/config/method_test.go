package config_test

import (
	"context"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/Goboolean/fetch-server/internal/adapter/broker"
	"github.com/Goboolean/fetch-server/internal/adapter/cache"
	"github.com/Goboolean/fetch-server/internal/adapter/meta"
	persistence_adapter "github.com/Goboolean/fetch-server/internal/adapter/persistence"
	"github.com/Goboolean/fetch-server/internal/adapter/transaction"
	"github.com/Goboolean/fetch-server/internal/adapter/websocket"
	"github.com/Goboolean/fetch-server/internal/domain/entity"
	"github.com/Goboolean/fetch-server/internal/domain/service/config"
	"github.com/Goboolean/fetch-server/internal/domain/service/persistence"
	"github.com/Goboolean/fetch-server/internal/domain/service/relayer"
	"github.com/Goboolean/fetch-server/internal/domain/service/transmission"
	"github.com/Goboolean/fetch-server/internal/infrastructure/ws/mock"
)



var instance *config.ConfigurationManager

func SetUp() {
	
	db           := persistence_adapter.NewMockAdapter()
	tx           := transaction.NewMock()
	meta         := meta.NewMockAdapter()
	ws := websocket.NewAdapter().(*websocket.Adapter)
	f := mock.New(context.Background(), time.Millisecond * 10, ws)

	if err := ws.RegisterFetcher(f); err != nil {
		panic(err)
	}

	relayer := relayer.New(db, tx, meta, ws)
	ws.RegisterReceiver(relayer)


	kafka := broker.NewMockAdapter()
	transmitter := transmission.New(kafka, relayer, transmission.Option{BatchSize: 2})

	cache    := cache.NewMockAdapter()
	persistenceManager := persistence.New(tx, db, cache, relayer, persistence.Option{BatchSize: 1})

	instance = config.New(meta, tx, relayer, persistenceManager, transmitter)
}



func TestMain(m *testing.M) {
	SetUp()
	code := m.Run()
	os.Exit(code)
}




func Test_StockConfiguration(t *testing.T) {


	var stockId = "stock.google.usa"
	t.Run("SetStockTransmittableTrue (case:false)", func(t *testing.T) {
		if err := instance.SetStockTransmittableTrue(context.Background(), stockId); err == nil {
			t.Errorf("SetStockTransmittableTrue() error = %v, expected error", err)
			return
		}
	})

	t.Run("SetStockStoreableTrue (case:false)", func(t *testing.T) {
		if err := instance.SetStockStoreableTrue(context.Background(), stockId); err == nil {
			t.Errorf("SetStockStoreableTrue() error = %v, expected error", err)
			return
		}
	})

	t.Run("SetStockRelayableTrue", func(t *testing.T) {
		if err := instance.SetStockRelayableTrue(context.Background(), stockId); err != nil {
			t.Errorf("SetStockRelayableTrue() error = %v", err)
			return
		}

		got, err := instance.GetStockConfiguration(context.Background(), stockId)
		if err != nil {
			t.Errorf("GetStockConfiguration() error = %v", err)
			return
		}

		want := entity.StockConfiguration{
			StockId: stockId,
			Relayable: true,
			Storeable: false,
			Transmittable: false,
		}

		if equals := reflect.DeepEqual(got, want); !equals {
			t.Errorf("GetStockConfiguration() = %v, want %v", got, want)
			return
		}
	})

	t.Run("SetStockTransmittableTrue (case:true)", func(t *testing.T) {
		if err := instance.SetStockTransmittableTrue(context.Background(), stockId); err != nil {
			t.Errorf("SetStockTransmittableTrue() error = %v", err)
			return
		}

		got, err := instance.GetStockConfiguration(context.Background(), stockId)
		if err != nil {
			t.Errorf("GetStockConfiguration() error = %v", err)
			return
		}

		want := entity.StockConfiguration{
			StockId: stockId,
			Relayable: true,
			Storeable: false,
			Transmittable: true,
		}

		if equals := reflect.DeepEqual(got, want); !equals {
			t.Errorf("GetStockConfiguration() = %v, want %v", got, want)
			return
		}
	})

	t.Run("SetStockStoreableTrue (case:true)", func(t *testing.T) {
		if err := instance.SetStockStoreableTrue(context.Background(), stockId); err != nil {
			t.Errorf("SetStockStoreableTrue() error = %v", err)
			return
		}

		got, err := instance.GetStockConfiguration(context.Background(), stockId)
		if err != nil {
			t.Errorf("GetStockConfiguration() error = %v", err)
			return
		}

		want := entity.StockConfiguration{
			StockId: stockId,
			Relayable: true,
			Storeable: true,
			Transmittable: true,
		}

		if equals := reflect.DeepEqual(got, want); !equals {
			t.Errorf("GetStockConfiguration() = %v, want %v", got, want)
			return
		}
	})

	t.Run("SetStockRelayableFalse", func(t *testing.T) {
		if err := instance.SetStockRelayableFalse(context.Background(), stockId); err != nil {
			t.Errorf("SetStockRelayableFalse() error = %v", err)
			return
		}

		got, err := instance.GetStockConfiguration(context.Background(), stockId)
		if err != nil {
			t.Errorf("GetStockConfiguration() error = %v", err)
			return
		}

		want := entity.StockConfiguration{
			StockId: stockId,
			Relayable: false,
			Storeable: false,
			Transmittable: false,
		}

		if equals := reflect.DeepEqual(got, want); !equals {
			t.Errorf("GetStockConfiguration() = %v, want %v", got, want)
			return
		}
	})

}