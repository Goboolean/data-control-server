package buycycle_test

import (
	"testing"
	"time"
)



func Test_SubscribeStockAggs(t *testing.T) {

	t.Skip("Skip this test, as buycycle server is not running.")

	const (
		symbol = "005930"
		falseSymbol = "000000"
	)

	var (
		countBeforeSubscription  int
		countAfterSubscription   int
		countAfterUnsubscription int
	)

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

		countBeforeSubscription = count
	
		time.Sleep(time.Second * 3/2)

		countAfterSubscription = count
		diff := countAfterSubscription - countBeforeSubscription

		if diff == 0 {
			t.Errorf(" received %d, want many", diff)
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
			t.Errorf("UnsubscribeStockAggs() = %v", err)
			return
		}

		time.Sleep(1 * time.Second)

		countAfterUnsubscription = count
		diff := countAfterUnsubscription - countAfterSubscription

		if diff != 0 {
			t.Errorf("UnsubscribeStockAggs() received %d, want 0", diff)
			return
		}
	})
}