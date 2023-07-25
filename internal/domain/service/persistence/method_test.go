package persistence_test

import (
	"testing"
)





func Test_SubscribeRelayer(t *testing.T) {

	stockId := "test"

	if err := instance.SubscribeRelayer(stockId); err != nil {
		t.Errorf("SubscribeRelayer() = %v", err)
		return
	}

	if storeable := instance.IsStockStoreable(stockId); !storeable {
		t.Errorf("IsStockStoreable() = %v, expected = true", storeable)
		return
	}

	if err := instance.UnsubscribeRelayer(stockId); err != nil {
		t.Errorf("UnsubscribeRelayer() = %v", err)
		return
	}

	if storeable := instance.IsStockStoreable(stockId); storeable {
		t.Errorf("IsStockStoreable() = %v, expected = false", storeable)
		return
	}

}



func Test_SynchronizationByDuration(t *testing.T) {

	stockId := "test"

	if err := instance.SubscribeRelayer(stockId); err != nil {
		t.Errorf("SubscribeRelayer() = %v", err)
		return
	}

}


func Test_SynchronizationByCount(t *testing.T) {

}