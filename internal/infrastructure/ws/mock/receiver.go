package mock

import "github.com/Goboolean/fetch-server/internal/infrastructure/ws"



// This class is not for product, but for testing.
type MockReceiver struct {
	f func()
}

func NewMockReceiver(f func()) *MockReceiver {
	return &MockReceiver{f: f}
}

func (m *MockReceiver) OnReceiveStockAggs(stock *ws.StockAggregate) error {
	m.f()
	return nil
}