package polygon_test

import (
	"testing"
	"time"

	"github.com/polygon-io/client-go/websocket/models"
)

var (
	counter  = 0
	stock    = "test"
)

type TestReceiver struct{}

func (r *TestReceiver) OnReceivePolygonStockAggs(models.EquityAgg) error {
	counter++
	return nil
}


func TestSubscribeStockAggs(t *testing.T) {

	if flag	:= isMarketOn(); !flag {
		t.Skip()
	}

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
