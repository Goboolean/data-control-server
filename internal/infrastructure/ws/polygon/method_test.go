package polygon_test

import (
	"testing"
	"time"

	"github.com/Goboolean/fetch-server.v1/internal/infrastructure/ws"
	"github.com/Goboolean/fetch-server.v1/internal/infrastructure/ws/mock"
)

var (
	count                = 0
	receiver ws.Receiver = mock.NewMockReceiver(func() {
		count++
	})
)

func Test_SubscribeStockAggs(t *testing.T) {

	const (
		symbol      = "AAPL"
		falseSymbol = "FALSE"
	)

	t.Skip("Skip this test, as polygon api key is expired.")

	t.Run("FalseSubscribe", func(t *testing.T) {
		if err := instance.SubscribeStockAggs(falseSymbol); err == nil {
			t.Errorf("SubscrbeStockAggs() = %v, want error", err)
			return
		}
	})

	t.Run("Subscribe", func(t *testing.T) {
		if err := instance.SubscribeStockAggs(symbol); err != nil {
			t.Errorf("SubscrbeStockAggs() = %v", err)
			return
		}

		countBeforeSubscription := count

		time.Sleep(time.Second * 2 / 3)

		countAfterSubscription := count
		diff := countAfterSubscription - countBeforeSubscription

		if diff == 0 {
			t.Errorf("SubscrbeStockAggs() received %d, want many", diff)
			return
		}
	})

	t.Run("SubscribeTwice", func(t *testing.T) {
		if err := instance.SubscribeStockAggs(symbol); err == nil {
			t.Errorf("SubscrbeStockAggs() = %v, want error", err)
			return
		}
	})

	t.Run("Unsubscribe", func(t *testing.T) {
		if err := instance.UnsubscribeStockAggs(symbol); err != nil {
			t.Errorf("UnsubscrbeStockAggs() = %v", err)
			return
		}

		countBeforeUnsubscription := count

		time.Sleep(time.Second * 2 / 3)

		countAfterUnsubscription := count
		diff := countAfterUnsubscription - countBeforeUnsubscription

		if diff != 0 {
			t.Errorf("UnsubscrbeStockAggs() received %d, want 0", diff)
			return
		}
	})

	t.Run("UnsubscribeTwice", func(t *testing.T) {
		if err := instance.UnsubscribeStockAggs(symbol); err == nil {
			t.Errorf("UnsubscrbeStockAggs() = %v, want error", err)
			return
		}
	})
}
