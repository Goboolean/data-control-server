package kis_test

import (
	"testing"
	"time"
)

var stockName = "DNASAAPL"

func Test_SubscribeStockAggs(t *testing.T) {

	if err := instance.SubscribeStockAggs(stockName); err != nil {
		t.Errorf("SubscrbeStockAggs() = %v", err)
		return
	}

	time.Sleep(1 * time.Second)

	if count == 0 {
		t.Errorf(" received %d, want many", count)
		return
	}
}


func Test_UnsubscribeStockAggs(t *testing.T) {

	
}
