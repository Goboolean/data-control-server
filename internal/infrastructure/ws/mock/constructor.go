package mock

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/Goboolean/fetch-server/internal/infrastructure/ws"
)



const platformName = "mock"

// MockFetcher is a mock implementation of Fetcher.
// It generates stock aggeregates data at an ramdom time with average of given duration.
// It does not fetch some data from external api, only generates.
// It is used for load test, performance test.
type MockFetcher struct {
	r ws.Receiver
	d time.Duration
	
	ctx context.Context
	cancel context.CancelFunc

	ch chan *ws.StockAggregate

	stocks map[string]*mockGenerater

	isClosed bool
}


func New(ctx context.Context, d time.Duration, r ws.Receiver) *MockFetcher {
	rand.Seed(time.Now().UnixNano())

	ctx, cancel := context.WithCancel(ctx)

	instance := &MockFetcher{
		d: d,
		r: r,
		ctx: ctx,
		cancel: cancel,
	}

	instance.ch = make(chan *ws.StockAggregate, 1000)
	instance.stocks = make(map[string]*mockGenerater)

	go instance.run()
	return instance
}


func (s *MockFetcher) PlatformName() string {
	return platformName
}


func (f *MockFetcher) run() {
	for {
		select {
		case <- f.ctx.Done():
			return
		case agg := <- f.ch:
			if err := f.r.OnReceiveStockAggs(agg); err != nil {
				log.Fatal(err)
			}
		}
	}
}


// Subscribing several topic at once is allowed, but atomicity is not guaranteed.
func (f *MockFetcher) SubscribeStockAggs(stocks ...string) error {

	for _, stock := range stocks {
		if _, ok := f.stocks[stock]; ok {
			return ErrTopicAlreadyExists
		}

		f.stocks[stock] = newMockGenerater(stock, f.ctx, f.ch, f.d)
	}
	return nil
}

// Unscribing several topic at once is allowed, but atomicity is not guaranteed.
func (f *MockFetcher) UnsubscribeStockAggs(stocks ...string) error {

	for _, stock := range stocks {
		if _, ok := f.stocks[stock]; !ok {
			return ErrTopicNotFound
		}

		f.stocks[stock].Close()
		delete(f.stocks, stock)
	}
	return nil
}


// Be sure to call Close() before the program ends.
func (f *MockFetcher) Close() error {
	// cancel() call will stop all the goroutines that generates data.
	f.cancel()
	close(f.ch)
	f.isClosed = true

	return nil
}


// MockFetcher does not explicitly connect to the server.
// Calling Ping() after Close() will cause error
func (f *MockFetcher) Ping() error {
	if flag := f.isClosed; flag {
		return ErrConnectionClosed
	}
	return nil
}