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

	instance = mock.New(context.Background(), 1*time.Millisecond, receiver)
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
	if err := instance.SubscribeStockAggs("AAPL"); err != nil {
		t.Errorf("SubscribeStockAggs() = %v", err)
	} 
}


func Test_UnsubscribeStockAggs(t *testing.T) {
	if err := instance.UnsubscribeStockAggs("AAPL"); err != nil {
		t.Errorf("UnsubscribeStockAggs() = %v", err)
	}
}