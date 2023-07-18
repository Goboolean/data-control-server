package websocket

import (
	"fmt"
	"sync"

	"github.com/Goboolean/fetch-server/internal/domain/port/in"
	"github.com/Goboolean/fetch-server/internal/infrastructure/ws"
)

// It is the only one entrypoint that stock data is sended to domain.
// It implements RelayerPort at domain port.
// It is compatible to any fetcher that implements Fetcher interface.
// Multiple fetchers can be used at the same time.
type StockWsAdapter struct {
	fetcher map[string]ws.Fetcher // stockid -> fetcher
	symbolToId map[string]string // symbol -> stockid

	port in.RelayerPort
}

var (
	instance *StockWsAdapter
	once     sync.Once
)

// There are two options to register fetcher:
// 1. compile time: use New()
// 2. runtime: use StockFetchAdapter.RegisterFetcher()
func New(port in.RelayerPort, fetchers ...ws.Fetcher) *StockWsAdapter {

	once.Do(func() {
		instance = &StockWsAdapter{
			fetcher: make(map[string]ws.Fetcher),
			symbolToId: make(map[string]string),
			port: port,
		}
	})

	for _, fetcher := range fetchers {
		name := fetcher.PlatformName()
		instance.fetcher[name] = fetcher
	}

	return instance
}


func (a *StockWsAdapter) RegisterFetcher(f ws.Fetcher) error {

	name := f.PlatformName()
	if _, ok := a.fetcher[name]; ok {
		return fmt.Errorf("fetcher %s is already registered", name)
	}

	a.fetcher[name] = f
	return nil
}


func (a *StockWsAdapter) UnregisterFetcher(f ws.Fetcher) error {

	name := f.PlatformName()
	if _, ok := a.fetcher[name]; !ok {
		return fmt.Errorf("fetcher %s is not registered", name)
	}

	delete(a.fetcher, name)
	return nil
}