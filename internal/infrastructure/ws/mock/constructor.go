package mock

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/Goboolean/fetch-server/internal/infrastructure/ws"
)

// MockFetcher is a mock implementation of Fetcher.
// It generates stock aggeregates data at an ramdom time with average of given duration.
type MockFetcher struct {
	r ws.Receiver
	d time.Duration
	
	ctx context.Context
	cancel context.CancelFunc

	ch chan *ws.StockAggregate

	stocks map[string]*mockGenerater
}


func New(ctx context.Context, d time.Duration, r ws.Receiver) *MockFetcher {
	rand.Seed(time.Now().UnixNano())

	ctx, cancel := context.WithCancel(ctx)

	stockChan := make(chan *ws.StockAggregate, 1000)

	go func() {
		for {
			select {
			case <- ctx.Done():
				return
			case agg := <- stockChan:
				if err := r.OnReceiveStockAggs(agg); err != nil {
					log.Fatal(err)
				}
			}
		}
	}()

	return &MockFetcher{
		d: d,
		r: r,
		ctx: ctx,
		cancel: cancel,
	}
}


func (f *MockFetcher) SubscribeStockAggs(stock string) error {
	if _, ok := f.stocks[stock]; ok {
		return errTopicAlreadyExists
	}

	f.stocks[stock] = newMockGenerater(stock, f.ctx, f.ch, f.d)
	return nil
}


func (f *MockFetcher) UnsubscribeStockAggs(stock string) error {
	if _, ok := f.stocks[stock]; !ok {
		return errTopicNotFound
	}

	f.stocks[stock].cancel()
	delete(f.stocks, stock)
	return nil
}


func (f *MockFetcher) Close() error {
	f.cancel()
	close(f.ch)

	return nil
}


func (f *MockFetcher) Ping() error {
	return nil
}