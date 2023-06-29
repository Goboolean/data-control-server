package polygon_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Goboolean/fetch-server/internal/infrastructure/polygon"
	"github.com/Goboolean/shared/pkg/resolver"
	"github.com/polygon-io/client-go/websocket/models"
)

var (
	instance *polygon.Subscriber
	counter  = 0
	stock    = "test"
)

type TestReceiver struct{}

func (r *TestReceiver) OnReceivePolygonStockAggs(models.EquityAgg) error {
	counter++
	return nil
}

func TestMain(m *testing.M) {

	instance = polygon.New(&resolver.Config{
		Host: os.Getenv("POLYGON_HOST"),
		Key:  os.Getenv("POLYGON_KEY"),
	}, &TestReceiver{})

	if !polygon.IsMarketOn() {
		fmt.Print("The foreign stock market is unavailable now")
		os.Exit(0)
	}

	code := m.Run()

	if err := instance.Close(); err != nil {
		panic(err)
	}

	os.Exit(code)
}

func TestSubscribeStockAggs(t *testing.T) {
	if err := instance.SubscribeStocksSecAggs(stock); err != nil {
		t.Errorf("SubscrbeStockAggs() = %v", err)
		return
	}

	time.Sleep(1 * time.Second)

	if counter == 0 {
		t.Errorf(" received %d, want many", counter)
		return
	}
}
