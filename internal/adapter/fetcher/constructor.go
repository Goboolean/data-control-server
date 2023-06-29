package fetcher

import (
	"sync"

	"github.com/Goboolean/fetch-server/internal/infrastructure/buycycle"
	"github.com/Goboolean/fetch-server/internal/infrastructure/polygon"
)

type StockFetchAdapter struct {
	b buycycle.Fetcher
	p polygon.Fetcher
}

var (
	instance *StockFetchAdapter
	once     sync.Once
)

func New(b buycycle.Fetcher, p polygon.Fetcher) *StockFetchAdapter {

	once.Do(func() {
		instance = &StockFetchAdapter{
			b: b,
			p: p,
		}
	})

	return instance
}
