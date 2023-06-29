package buycycle_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Goboolean/fetch-server/internal/infrastructure/buycycle"
	"github.com/Goboolean/shared/pkg/resolver"
)

var (
	instance *buycycle.Subscriber
	counter  = 0
	stock    = "test"
)

type TestReceiver struct{}

func (r *TestReceiver) OnReceiveBuycycleStockAggs(buycycle.StockAggregate) error {
	counter++
	return nil
}

func TestMain(m *testing.M) {

	instance = buycycle.New(&resolver.Config{
		Host: os.Getenv("BUYCYCLE_HOST"),
		Port: os.Getenv("BUYCYCLE_PORT"),
		Path: os.Getenv("BUYCYCLE_PATH"),
	}, &TestReceiver{})

	if !buycycle.IsMarketOn() {
		fmt.Print("The domestic stock market is unavailable now")
		os.Exit(0)
	}

	code := m.Run()

	if err := instance.Close(); err != nil {
		panic(err)
	}

	os.Exit(code)
}

func TestSubscribeStockAggs(t *testing.T) {
	if err := instance.SubscribeStockAggs(stock); err != nil {
		t.Errorf("SubscrbeStockAggs() = %v", err)
		return
	}

	time.Sleep(1 * time.Second)

	if counter == 0 {
		t.Errorf(" received %d, want many", counter)
		return
	}
}
