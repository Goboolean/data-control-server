package transmission_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Goboolean/fetch-server/internal/adapter/broker"
	"github.com/Goboolean/fetch-server/internal/adapter/meta"
	"github.com/Goboolean/fetch-server/internal/adapter/persistence"
	"github.com/Goboolean/fetch-server/internal/adapter/transaction"
	"github.com/Goboolean/fetch-server/internal/adapter/websocket"
	"github.com/Goboolean/fetch-server/internal/domain/port/out"
	"github.com/Goboolean/fetch-server/internal/domain/service/relay"
	"github.com/Goboolean/fetch-server/internal/domain/service/transmission"
	"github.com/Goboolean/fetch-server/internal/infrastructure/ws/mock"
)

var (
	instance       *transmission.Manager
	kafka          out.TransmissionPort
	relayer *relay.Manager
)

func MockRelayer() *relay.Manager {

	var (
		db   = persistence.NewMockAdapter()
		tx   = transaction.NewMock()
		meta = meta.NewMockAdapter()
		ws   = websocket.NewMockAdapter().(*websocket.MockAdapter)
		f    = mock.New(time.Millisecond*10, ws)
	)

	if err := ws.RegisterFetcher(f); err != nil {
		panic(err)
	}

	instance, err := relay.New(db, tx, meta, ws)
	if err != nil {
		panic(err)
	}
	ws.RegisterReceiver(instance)

	return instance
}

func SetUp() {
	var err error

	var stockId = "stock.google.usa"

	relayer = MockRelayer()
	kafka = broker.NewMockAdapter()

	instance, err = transmission.New(kafka, relayer, transmission.Option{BatchSize: 2})
	if err != nil {
		panic(err)
	}

	if err := relayer.FetchStock(context.Background(), stockId); err != nil {
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

func Test_Manager(t *testing.T) {

	var stockId = "stock.google.usa"

	var count = 0

	t.Run("SubscribeRelayer", func(t *testing.T) {
		if err := instance.SubscribeRelayer(context.Background(), stockId); err != nil {
			t.Errorf("SubscribeRelayer() = %v", err)
			return
		}

		time.Sleep(time.Millisecond * 100)

		sended, err := kafka.(*broker.MockAdapter).GetTransmittedStockCount(stockId)
		if err != nil {
			t.Errorf("error getting transmitted stock count: %s", err)
			return
		}

		if sended == 0 {
			t.Errorf("received %d, want many", sended)
			return
		}

		count = sended
	})

	t.Run("UnsubscribeRelayer", func(t *testing.T) {
		if err := instance.UnsubscribeRelayer(stockId); err != nil {
			t.Errorf("UnsubscribeRelayer() = %v", err)
			return
		}

		time.Sleep(time.Millisecond * 100)

		sended, err := kafka.(*broker.MockAdapter).GetTransmittedStockCount(stockId)
		if err != nil {
			t.Errorf("error getting transmitted stock count: %s", err)
			return
		}

		if sended != count {
			t.Errorf("received %d, want %d", sended, count)
			return
		}
	})

	t.Run("IsStockTransmittable", func(t *testing.T) {
		if flag := instance.IsStockTransmittable(stockId); flag {
			t.Errorf("IsStockTransmittable() = %v, expected = false", flag)
			return
		}

		if err := instance.SubscribeRelayer(context.Background(), stockId); err != nil {
			t.Errorf("SubscribeRelayer() = %v", err)
			return
		}

		if err := instance.UnsubscribeRelayer(stockId); err != nil {
			t.Errorf("UnsubscribeRelayer() = %v", err)
			return
		}

		if flag := instance.IsStockTransmittable(stockId); flag {
			t.Errorf("IsStockTransmittable() = %v, expected = %v", flag, false)
			return
		}
	})

	t.Run("IsStockTransmittable (whenSetTransmission)", func(t *testing.T) {
		if flag := instance.IsStockTransmittable(stockId); flag {
			t.Errorf("IsStockTransmittable() = %v, expected = false", flag)
			return
		}

		if err := relayer.StopFetchingStock(context.Background(), stockId); err != nil {
			t.Errorf("StopFetchingStock() = %v", err)
			return
		}

		time.Sleep(time.Millisecond * 100)

		if flag := instance.IsStockTransmittable(stockId); flag {
			t.Errorf("IsStockTransmittable() = %v, expected = %v", flag, false)
			return
		}
	})

}
