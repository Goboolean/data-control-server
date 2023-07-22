package websocket

import (
	"context"
	"fmt"

	"github.com/Goboolean/fetch-server/internal/domain/entity"
	"github.com/Goboolean/fetch-server/internal/domain/port/in"
	"github.com/Goboolean/fetch-server/internal/domain/port/out"
	"github.com/Goboolean/fetch-server/internal/infrastructure/prometheus"
	"github.com/Goboolean/fetch-server/internal/infrastructure/ws"
)

// It is the only one entrypoint that stock data is sended to domain.
// It implements RelayerPort at domain port.
// It is compatible to any fetcher that implements Fetcher interface.
// Multiple fetchers can be used at the same time.
type Adapter struct {
	fetcher map[string]ws.Fetcher // stockid -> fetcher
	symbolToId map[string]string // symbol -> stockid

	port in.RelayerPort
}


// There are two options to register fetcher:
// 1. compile time: use New()
// 2. runtime: use StockFetchAdapter.RegisterFetcher()
func NewAdapter(port in.RelayerPort, fetchers ...ws.Fetcher) out.RelayerPort {

	instance := &Adapter{
		fetcher: make(map[string]ws.Fetcher),
		symbolToId: make(map[string]string),
		port: port,
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
		return fmt.Errorf("fetcher %s is already registered", name)
	}

	a.fetcher[name] = f
	return nil
}


func (a *Adapter) UnregisterFetcher(f ws.Fetcher) error {

	name := f.PlatformName()
	if _, ok := a.fetcher[name]; !ok {
		return fmt.Errorf("fetcher %s is not registered", name)
	}

	delete(a.fetcher, name)
	return nil
}



func (s *Adapter) toDomainEntity(agg *ws.StockAggregate) (*entity.StockAggregateForm, error) {
	stockId, ok := s.symbolToId[agg.Symbol]
	if !ok {
		return nil, ErrStockNotFound
	}

	return &entity.StockAggregateForm{
		StockAggsMeta: entity.StockAggsMeta{
			StockID: stockId,
		},
		StockAggregate: entity.StockAggregate{
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


func (s *Adapter) OnReceiveStockAggs(agg *ws.StockAggregate) error {
	prometheus.DomesticStockCounter.Inc()

	data, err := s.toDomainEntity(agg)
	if err != nil {
		return err
	}

	return s.port.PlaceStockFormBatch([]*entity.StockAggregateForm{data})
}

func (s *Adapter) OnReceiveStockAggsBatch(aggs []*ws.StockAggregate) error {
	prometheus.DomesticStockCounter.Add(float64(len(aggs)))

	batch := make([]*entity.StockAggregateForm, len(aggs))

	for _, agg := range aggs {
		data, err := s.toDomainEntity(agg)
		if err != nil {
			return err
		}

		batch = append(batch, data)
	}

	return s.port.PlaceStockFormBatch(batch)
}


func (s *Adapter) FetchStock(ctx context.Context, stockId string, stockMeta entity.StockAggsMeta) error {
	platform := stockMeta.Platform
	fetcher, ok := s.fetcher[platform]

	if !ok {
		return fmt.Errorf("fetcher %s is not registered", platform)
	}

	return fetcher.SubscribeStockAggs(stockId)
}


func (s *Adapter) StopFetchingStock(ctx context.Context, stockId string) error {
	platform, ok := s.symbolToId[stockId]
	if !ok {
		return ErrStockNotFound
	}

	fetcher, ok := s.fetcher[platform]
	if !ok {
		return fmt.Errorf("fetcher %s is not registered", platform)
	}

	return fetcher.UnsubscribeStockAggs(stockId)
}

