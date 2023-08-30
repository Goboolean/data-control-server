package websocket

import (
	"context"

	"github.com/Goboolean/fetch-server.v1/internal/domain/port/in"
	"github.com/Goboolean/fetch-server.v1/internal/domain/vo"
	"github.com/Goboolean/fetch-server.v1/internal/infrastructure/ws"
	"github.com/Goboolean/fetch-server.v1/internal/util/prometheus"
)

// It is the only one entrypoint that stock data is sended to domain.
// It implements RelayerPort at domain port.
// It is compatible to any fetcher that implements Fetcher interface.
// Multiple fetchers can be used at the same time.
type Adapter struct {
	fetcher      map[string]ws.Fetcher // stockid -> fetcher
	symbolToId   map[string]string     // symbol -> stockid
	idToPlatform map[string]string     // stockid -> platform
	idToSymnbol  map[string]string     // stockid -> symbol

	port in.RelayerPort
}

// There are two options to register fetcher:
// 1. compile time: use New()
// 2. runtime: use StockFetchAdapter.RegisterFetcher()
func NewAdapter(fetchers ...ws.Fetcher) *Adapter {

	instance := &Adapter{
		fetcher:      make(map[string]ws.Fetcher),
		symbolToId:   make(map[string]string),
		idToPlatform: make(map[string]string),
		idToSymnbol:  make(map[string]string),
	}

	for _, fetcher := range fetchers {
		name := fetcher.PlatformName()
		instance.fetcher[name] = fetcher
	}

	return instance
}

func (a *Adapter) RegisterFetcher(f ws.Fetcher) error {

	name := f.PlatformName()
	if _, ok := a.fetcher[name]; ok {
		return ErrFetcherAlreadyRegistered
	}

	a.fetcher[name] = f
	return nil
}

func (a *Adapter) UnregisterFetcher(f ws.Fetcher) error {

	name := f.PlatformName()
	if _, ok := a.fetcher[name]; !ok {
		return ErrFetcherNotRegistered
	}

	delete(a.fetcher, name)
	return nil
}

func (a *Adapter) RegisterReceiver(port in.RelayerPort) {
	a.port = port
}

func (s *Adapter) toDomainEntity(agg *ws.StockAggregate) (*vo.StockAggregateForm, error) {
	stockId, ok := s.symbolToId[agg.Symbol]
	if !ok {
		return nil, ErrSymbolUnrecognized
	}

	return &vo.StockAggregateForm{
		StockAggsMeta: vo.StockAggsMeta{
			StockID: stockId,
		},
		StockAggregate: vo.StockAggregate{
			Min:    agg.Min,
			Max:    agg.Max,
			Open:   agg.Start,
			Closed: agg.End,
			Time:   agg.StartTime,
		},
	}, nil
}

func (s *Adapter) OnReceiveStockAggs(agg *ws.StockAggregate) error {

	data, err := s.toDomainEntity(agg)
	if err != nil {
		return err
	}

	s.port.PlaceStockFormBatch([]*vo.StockAggregateForm{data})

	prometheus.FetchCounter.Inc()
	return nil
}

func (s *Adapter) OnReceiveStockAggsBatch(aggs []*ws.StockAggregate) error {
	batch := make([]*vo.StockAggregateForm, len(aggs))

	for _, agg := range aggs {
		data, err := s.toDomainEntity(agg)
		if err != nil {
			return err
		}

		batch = append(batch, data)
	}

	s.port.PlaceStockFormBatch(batch)

	prometheus.FetchCounter.Add(float64(len(aggs)))
	return nil
}

func (s *Adapter) FetchStock(ctx context.Context, stockId string, platform string, symbol string) error {

	fetcher, ok := s.fetcher[platform]

	if !ok {
		return ErrFetcherNotRegistered
	}

	if err := fetcher.SubscribeStockAggs(symbol); err != nil {
		return err
	}

	s.idToPlatform[stockId] = platform
	s.symbolToId[symbol] = stockId
	s.idToSymnbol[stockId] = symbol

	return nil
}

func (s *Adapter) StopFetchingStock(ctx context.Context, stockId string) error {
	platform, ok := s.idToPlatform[stockId]
	if !ok {
		return ErrPlatformNotFoundByStockId
	}

	fetcher, ok := s.fetcher[platform]
	if !ok {
		return ErrFetcherNotFoundByPlatformName
	}

	symbol, ok := s.idToSymnbol[stockId]
	if !ok {
		return ErrSymbolNotFoundByStockId
	}

	return fetcher.UnsubscribeStockAggs(symbol)
}
