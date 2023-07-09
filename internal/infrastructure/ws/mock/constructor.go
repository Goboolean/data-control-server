package mock

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/Goboolean/fetch-server/internal/infrastructure/ws"
	"github.com/Goboolean/fetch-server/internal/util/math"
)

// MockFetcher is a mock implementation of Fetcher.
// It generates stock aggeregates data at an ramdom time with average of given duration.
type MockFetcher struct {
	r ws.Receiver
	d time.Duration
	
	ctx context.Context
	cancel context.CancelFunc
}



func New(ctx context.Context, d time.Duration, r ws.Receiver) *MockFetcher {
	rand.Seed(time.Now().UnixNano())

	ctx, cancel := context.WithCancel(ctx)

	return &MockFetcher{
		d: d,
		r: r,
		ctx: ctx,
		cancel: cancel,
	}
}


func (f *MockFetcher) SubscribeStockAggs(stock string) error {
	return nil
}



var (
	lastTime time.Time = time.Now()
	lastPrice float64 = 1000
)

func (f *MockFetcher) generateRandomStockAggs() *ws.StockAggregate {

	curTime := time.Now()
	curPrice := lastPrice * (rand.Float64() * 0.2 + 0.9)

	stockAggs := &ws.StockAggregate{
		StartTime: lastTime.UnixNano(),
		EndTime: curTime.UnixNano(),
		Average: (lastPrice + curPrice) / 2,
		Min: math.MinFloat(lastPrice, curPrice),
		Max: math.MaxFloat(lastPrice, curPrice),
		Start: lastPrice,
		End: curPrice,
	}

	lastTime = curTime
	lastPrice = curPrice

	return stockAggs
}


func (f *MockFetcher) Run() {

	go func() {
		for {
			newDuration := time.Duration(rand.Int63n(2 * int64(f.d)))

			select {
			
			case <- f.ctx.Done():
				return

			case <- time.After(newDuration):
				stockAggs := f.generateRandomStockAggs()
				if err := f.r.OnReceiveStockAggs(stockAggs); err != nil {
					log.Fatal(err)
				}
			}
		}
	}()

	return
}


func (f *MockFetcher) Close() error {
	f.cancel()
	return nil
}


func (f *MockFetcher) Ping() error {
	return nil
}