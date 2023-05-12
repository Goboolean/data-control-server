package mocks

import (
	"context"
	"sync"
	"time"

	"github.com/Goboolean/stock-fetch-server/internal/infrastructure/buycycle"
)


type MockSubscriber struct {

	ctx context.Context
	cancel context.CancelFunc
	r buycycle.Receiver
}

var (
	instance *MockSubscriber
	once sync.Once
)

// Create arbitary data with a period of given duration, and call r.OnReceiveBuycycle() method
func New(duration time.Duration, r buycycle.Receiver) *MockSubscriber {

	once.Do(func() {
		ctx, cancel := context.WithCancel(context.Background())

		instance = &MockSubscriber{
			ctx: ctx,
			cancel: cancel,
			r: r,
		}
	})

	return instance
}



func (s *MockSubscriber) Close() error {
	if err := s.Close(); err != nil {
		return err
	}

	s.cancel()
	return nil
}

