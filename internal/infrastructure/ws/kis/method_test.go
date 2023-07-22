package kis_test

import (
	"testing"
	"time"
)

var stockName = "DNASAAPL"

func Test_SubscribeStockAggs(t *testing.T) {
	/*
		if flag := isMarketOn(); flag {
			t.Skip()
			return
		}
	*/
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

// UnsubscribeStockAggs is not implemented yet, so this test is skipped.
func Test_UnsubscribeStockAggs(t *testing.T) {
	t.Skip()
}
