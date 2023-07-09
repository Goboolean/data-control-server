package mock_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Goboolean/fetch-server/internal/infrastructure/ws"
	"github.com/Goboolean/fetch-server/internal/infrastructure/ws/mock"
)

var instance ws.Fetcher
	
var (
	receiver ws.Receiver 
	count int = 0
)



func SetupMock() {
	receiver = mock.NewMockReceiver(func() {
		count++
	})

	instance = mock.New(context.Background(), 10*time.Millisecond, receiver)
}


func TeardownMock() {
	instance.Close()
}


func TestMain(m *testing.M) {

	SetupMock()
	code := m.Run()
	TeardownMock()

	os.Exit(code)
}


func Test_Constructor(t *testing.T) {

	if err := instance.Ping(); err != nil {
		t.Errorf("Ping() = %v", err)
	}
}




func Test_SubscribeStockAggs(t *testing.T) {
	// I thought count should be in the interval of [5,20],
  // since the data is generated every 10 ms in average.
  // But it seems that error may occur in some case.

	if err := instance.SubscribeStockAggs("test"); err != nil {
		t.Errorf("SubscribeStockAggs() = %v", err)
	}

	if err := instance.SubscribeStockAggs("test"); err == nil {
		t.Errorf("SubscribeStockAggs() = %v, want %v", err, mock.ErrTopicAlreadyExists)
	}

	time.Sleep(100 * time.Millisecond)

	if !(5 <= count) {
		t.Errorf("count = %v, should be at least 5", count)
	}
}


func Test_UnsubscribeStockAggs(t *testing.T) {

	if err := instance.UnsubscribeStockAggs("unsubscribed"); err == nil {
		t.Errorf("UnsubscribeStockAggs() = %v, want %v", err, mock.ErrTopicNotFound)
	}

	if err := instance.SubscribeStockAggs("test"); err != nil {
		t.Errorf("SubscribeStockAggs() = %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	if err := instance.UnsubscribeStockAggs("test"); err != nil {
		t.Errorf("UnsubscribeStockAggs() = %v", err)
	}

	lastCount := count

	time.Sleep(100 * time.Millisecond)

	if lastCount != count {
		t.Errorf("lastCount = %v, count = %v", lastCount, count)
	}
}