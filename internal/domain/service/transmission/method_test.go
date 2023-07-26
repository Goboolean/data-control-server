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
	"github.com/Goboolean/fetch-server/internal/domain/service/relayer"
	"github.com/Goboolean/fetch-server/internal/domain/service/transmission"
	"github.com/Goboolean/fetch-server/internal/infrastructure/ws/mock"
)



var (
	instance *transmission.Transmitter
	kafka out.TransmissionPort
	relayerManager *relayer.RelayerManager
)

func MockRelayer() *relayer.RelayerManager {

	var (
		db           = persistence.NewMockAdapter()
		tx           = transaction.NewMock()
		meta         = meta.NewMockAdapter()
		ws = websocket.NewAdapter()
		f = mock.New(context.Background(), time.Millisecond * 10, ws)
	)

	if err := ws.RegisterFetcher(f); err != nil {
		panic(err)
	}

	instance := relayer.New(db, tx, meta, ws)
	ws.RegisterReceiver(instance)

	return instance
}


func SetUp() {

	var stockId = "stock.google.usa"
				
	relayerManager = MockRelayer()
	kafka = broker.NewMockAdapter()
	instance = transmission.New(context.Background(), kafka, relayerManager, transmission.Option{BatchSize: 2})

	if err := relayerManager.FetchStock(context.Background(), stockId); err != nil {
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



func Test_Transmitter(t *testing.T) {

	var stockId = "stock.google.usa"

	var count = 0

	
	t.Run("SubscribeRelayer", func(t *testing.T) {
		if err := instance.SubscribeRelayer(context.Background(), stockId); err != nil {
			t.Errorf("SubscribeRelayer() = %v", err)
			return
		}

		time.Sleep(time.Millisecond * 100)

		sended, err := kafka.(*broker.MockAdapter).GetTransmittedStockCount(stockId); 
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

		sended, err := kafka.(*broker.MockAdapter).GetTransmittedStockCount(stockId);
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

		if err := relayerManager.StopFetchingStock(context.Background(), stockId); err != nil {
			t.Errorf("StopFetchingStock() = %v", err)
			return
		}

		time.Sleep(time.Millisecond * 100)

		if flag := instance.IsStockTransmittable(stockId); !flag {
			t.Errorf("IsStockTransmittable() = %v, expected = %v", flag, true)
			return
		}
	})

}

