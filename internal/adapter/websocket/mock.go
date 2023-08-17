package websocket

import (
	"context"

	"github.com/Goboolean/fetch-server/internal/domain/vo"
	"github.com/Goboolean/fetch-server/internal/domain/port/in"
	"github.com/Goboolean/fetch-server/internal/domain/port/out"
	"github.com/Goboolean/fetch-server/internal/infrastructure/ws"
)

// It is the only one entrypoint that stock data is sended to domain.
// It implements RelayerPort at domain port.
// It is compatible to any fetcher that implements Fetcher interface.
// Multiple fetchers can be used at the same time.
type MockAdapter struct {
	fetcher map[string]ws.Fetcher // stockid -> fetcher
	symbolToId map[string]string // symbol -> stockid
	idToPlatform map[string]string // stockid -> platform
	idToSymnbol map[string]string // stockid -> symbol

	port in.RelayerPort
}


// There are two options to register fetcher:
// 1. compile time: use New()
// 2. runtime: use StockFetchAdapter.RegisterFetcher()
func NewMockAdapter() out.RelayerPort {

	instance := &MockAdapter{
		fetcher: make(map[string]ws.Fetcher),
		symbolToId: make(map[string]string),
		idToPlatform: make(map[string]string),
		idToSymnbol: make(map[string]string),
	}

	return instance
}


func (a *MockAdapter) RegisterFetcher(f ws.Fetcher) error {

	name := f.PlatformName()
	if _, ok := a.fetcher[name]; ok {
		return ErrFetcherAlreadyRegistered
	}

	a.fetcher[name] = f
	return nil
}


func (a *MockAdapter) UnregisterFetcher(f ws.Fetcher) error {

	name := f.PlatformName()
	if _, ok := a.fetcher[name]; !ok {
		return ErrFetcherNotRegistered
	}

	delete(a.fetcher, name)
	return nil
}


func (a *MockAdapter) RegisterReceiver(port in.RelayerPort) {
	a.port = port
}


func (s *MockAdapter) toDomainEntity(agg *ws.StockAggregate) (*vo.StockAggregateForm, error) {
	stockId, ok := s.symbolToId[agg.Symbol]
	if !ok {
		return nil, ErrSymbolUnrecognized
	}

	return &vo.StockAggregateForm{
		StockAggsMeta: vo.StockAggsMeta{
			StockID: stockId,
		},
		StockAggregate: vo.StockAggregate{
			Average: agg.Average,
			Min: agg.Min,
			Max: agg.Max,
			Start: agg.Start,
			End: agg.End,
			StartTime: agg.StartTime,
			EndTime: agg.EndTime,
		},
	}, nil
}


func (s *MockAdapter) OnReceiveStockAggs(agg *ws.StockAggregate) error {

	data, err := s.toDomainEntity(agg)
	if err != nil {
		return err
	}

	s.port.PlaceStockFormBatch([]*vo.StockAggregateForm{data})
	return nil
}

func (s *MockAdapter) OnReceiveStockAggsBatch(aggs []*ws.StockAggregate) error {
	batch := make([]*vo.StockAggregateForm, len(aggs))

	for _, agg := range aggs {
		data, err := s.toDomainEntity(agg)
		if err != nil {
			return err
		}

		batch = append(batch, data)
	}

	s.port.PlaceStockFormBatch(batch)
	return nil
}


func (s *MockAdapter) FetchStock(ctx context.Context, stockId string, platform string, symbol string) error {

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



func (s *MockAdapter) StopFetchingStock(ctx context.Context, stockId string) error {
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

