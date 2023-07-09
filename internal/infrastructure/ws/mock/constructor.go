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

	go func() {
		for {
			select {
			case <- ctx.Done():
				return
			case agg := <- instance.ch:
				if err := r.OnReceiveStockAggs(agg); err != nil {
					log.Fatal(err)
				}
			}
		}
	}()

	return instance
}


// Subscribing several topic at once is allowed, but atomicity is not guaranteed.
func (f *MockFetcher) SubscribeStockAggs(stocks ...string) error {

	for _, stock := range stocks {
		if _, ok := f.stocks[stock]; ok {
			return errTopicAlreadyExists
		}

		f.stocks[stock] = newMockGenerater(stock, f.ctx, f.ch, f.d)
	}
	return nil
}

// Unscribing several topic at once is allowed, but atomicity is not guaranteed.
func (f *MockFetcher) UnsubscribeStockAggs(stocks ...string) error {

	for _, stock := range stocks {
		if _, ok := f.stocks[stock]; !ok {
			return errTopicNotFound
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
		return errConnectionClosed
	}
	return nil
}