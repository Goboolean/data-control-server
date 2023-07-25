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

	var (		
		relayer = MockRelayer()
		stockId = "stock.google.usa"
	)

	kafka = broker.NewMockAdapter()
	instance = transmission.New(context.Background(), kafka, relayer, transmission.Option{BatchSize: 1})

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



func Test_Transmitter(t *testing.T) {

	var stockId = "stock.google.usa"

	if flag := instance.IsStockTransmittable(stockId); flag {
		t.Errorf("stockId %s should not be transmittable", stockId)
		return
	}

	if err := instance.SubscribeRelayer(context.Background(), stockId); err != nil {
		t.Errorf("error subscribing relayer: %s", err)
		return
	}

	time.Sleep(time.Millisecond * 100)

	if flag := instance.IsStockTransmittable(stockId); !flag {
		t.Errorf("stockId %s should be transmittable", stockId)
		return
	}

	sended, err := kafka.(*broker.MockAdapter).GetTransmittedStockCount(stockId); 
	if err != nil {
		t.Errorf("error getting transmitted stock count: %s", err)
		return
	}
	
	if sended == 0 {
		t.Errorf("received %d, want many", sended)
		return
	}

	if err := instance.UnsubscribeRelayer(stockId); err != nil {
		t.Errorf("error unsubscribing relayer: %s", err)
		return
	}

	if flag := instance.IsStockTransmittable(stockId); flag {
		t.Errorf("stockId %s should not be transmittable", stockId)
		return
	}

}

