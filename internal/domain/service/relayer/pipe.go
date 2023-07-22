package relayer

import (
	"context"
	"fmt"

	"github.com/Goboolean/fetch-server/internal/domain/entity"
)

const BUFFER_SIZE = 100

type pipe struct {
	filteredChan      chan *entity.StockAggregateForm
	classifiedChanMap map[string]chan *entity.StockAggregate

	startPoint chan *entity.StockAggregateForm
	endPoint   map[string]chan []*entity.StockAggregate
}

func (p *pipe) filterBadTick(in <-chan *entity.StockAggregateForm, out chan<- *entity.StockAggregateForm) {
	for stock := range in {
		out <- stock
	}
}

func (p *pipe) classifyStock(in <-chan *entity.StockAggregateForm, out map[string]chan *entity.StockAggregate) {
	for stock := range in {
		out[stock.StockID] <- &stock.StockAggregate
	}
}

func (p *pipe) bindBatch(in <-chan *entity.StockAggregate, out chan<- []*entity.StockAggregate) {
	batch := make([]*entity.StockAggregate, 0)

	for stock := range in {
		batch = append(batch, stock)
		if len(batch) == BUFFER_SIZE {
			out <- batch
			batch = batch[:0]
		}
	}
}

func newPipe(st chan *entity.StockAggregateForm, ed map[string]chan []*entity.StockAggregate) *pipe {
	instance := &pipe{
		filteredChan:      make(chan *entity.StockAggregateForm),
		classifiedChanMap: make(map[string]chan *entity.StockAggregate),
		startPoint:        st,
		endPoint:          ed,
	}

	return instance
}

// Run as goroutine, and control lifeccle with ctx.
func (p *pipe) ExecPipe(ctx context.Context) {

	for stock := range p.endPoint {
		p.classifiedChanMap[stock] = make(chan *entity.StockAggregate)
	}

	go p.filterBadTick(p.startPoint, p.filteredChan)
	go p.classifyStock(p.filteredChan, p.classifiedChanMap)

	// Use AddNewPipe() to create filter fot channels to met endpoint

	defer func() {
		for stock := range p.endPoint {
			close(p.classifiedChanMap[stock])
			delete(p.classifiedChanMap, stock)
		}
		close(p.filteredChan)
	}()

	<-ctx.Done()
}

func (p *pipe) AddNewPipe(stock string) {
	p.classifiedChanMap[stock] = make(chan *entity.StockAggregate)
	go p.bindBatch(p.classifiedChanMap[stock], p.endPoint[stock])
}

func (p *pipe) PlaceOnStartPoint(data *entity.StockAggregateForm) {
	p.startPoint <- data
}

func (p *pipe) GetEndpointChannel(stock string) (chan []*entity.StockAggregate, error) {
	chans, ok := p.endPoint[stock]
	if !ok {
		return nil, fmt.Errorf("stock do not exist on channel")
	}

	return chans, nil
}
