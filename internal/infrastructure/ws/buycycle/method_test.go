package buycycle_test

import (
	"testing"
	"time"

	"github.com/Goboolean/fetch-server/internal/infrastructure/ws/buycycle"
)

var (
	counter  = 0
	stock    = "test"
)

type TestReceiver struct{}

func (r *TestReceiver) OnReceiveBuycycleStockAggs(buycycle.StockAggregate) error {
	counter++
	return nil
}



func TestSubscribeStockAggs(t *testing.T) {

	if flag := isMarketOn(); flag {
		t.Skip()
	}

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
