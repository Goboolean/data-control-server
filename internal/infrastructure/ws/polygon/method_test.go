package polygon_test

import (
	"testing"
	"time"

	"github.com/Goboolean/fetch-server/internal/infrastructure/ws"
	"github.com/Goboolean/fetch-server/internal/infrastructure/ws/mock"
)



var (
	count = 0
	receiver ws.Receiver = mock.NewMockReceiver(func() {
		count++
	})
)


func Test_SubscribeStockAggs(t *testing.T) {

	const stock1 = "stock1"
	const stock2 = "stock2"

	if flag	:= isMarketOn(); !flag {
		t.Skip()
	}

	if err := instance.SubscribeStockAggs(stock1, stock2); err != nil {
		t.Errorf("SubscrbeStockAggs() = %v", err)
		return
	}

	time.Sleep(2 * time.Second)

	if count == 0 {
		t.Errorf("SubscrbeStockAggs() received %d, want many", count)
		return
	}
}


func Test_UnsubscribeStockAggs(t *testing.T) {

	const stock1 = "stock1"
	const stock2 = "stock2"

	if flag	:= isMarketOn(); !flag {
		t.Skip()
	}

	if err := instance.SubscribeStockAggs(stock1, stock2); err != nil {
		t.Errorf("SubscrbeStockAggs() = %v", err)
		return
	}

	if err := instance.UnsubscribeStockAggs(stock2, stock1); err != nil {
		t.Errorf("UnsubscrbeStockAggs() = %v", err)
		return
	}	

	time.Sleep(2 * time.Second)

	if count != 0 {
		t.Errorf("UnsubscrbeStockAggs() received %d, want 0", count)
		return
	}
}