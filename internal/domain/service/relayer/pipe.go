package relayer

import (
	"context"
	"fmt"

	"github.com/Goboolean/stock-fetch-server/internal/domain/value"
)


const BUFFER_SIZE = 100

type pipe struct {
	filteredChan chan value.StockAggregateForm
	classifiedChanMap map[string] chan value.StockAggregate

	startPoint chan value.StockAggregateForm
	endPoint map[string] chan []value.StockAggregate
}



func (p *pipe) filterBadTick(in <-chan value.StockAggregateForm, out chan<- value.StockAggregateForm) {
	for stock := range in {
		out <- stock
	}
}

func (p *pipe) classifyStock(in <-chan value.StockAggregateForm, out map[string] chan value.StockAggregate) {
	for stock := range in {
		out[stock.StockID] <- stock.StockAggregate
	}
}

func (p *pipe) bindBatch(in <-chan value.StockAggregate, out chan<- []value.StockAggregate) {
	batch := make([]value.StockAggregate, 0)

	for stock := range in {
		batch = append(batch, stock)
		if len(batch) == BUFFER_SIZE {
			out <- batch
			batch = batch[:0]
		}
	}
}

func newPipe(st chan value.StockAggregateForm, ed map[string] chan []value.StockAggregate) *pipe {
	instance := &pipe{
		filteredChan: make(chan value.StockAggregateForm),
		classifiedChanMap: make(map[string] chan value.StockAggregate),
		startPoint: st,
		endPoint: ed,
	}

	return instance
}


// Run as goroutine, and control lifeccle with ctx.
func (p *pipe) ExecPipe(ctx context.Context) {

	for stock := range p.endPoint {
		p.classifiedChanMap[stock] = make(chan value.StockAggregate)
	}

	go p.filterBadTick(p.startPoint, p.filteredChan)
	go p.classifyStock(p.filteredChan, p.classifiedChanMap)

	// Use AddNewPipe() to create filter fot channels to met endpoint

	select {
	case <- ctx.Done():

		for stock := range p.endPoint {
			close(p.classifiedChanMap[stock])
			delete(p.classifiedChanMap, stock)
		}

		close(p.filteredChan)
	}
}



func (p *pipe) AddNewPipe(stock string) {
	p.classifiedChanMap[stock] = make(chan value.StockAggregate)
	go p.bindBatch(p.classifiedChanMap[stock], p.endPoint[stock])
}

func (p *pipe) PlaceOnStartPoint(data value.StockAggregateForm) {
	p.startPoint <- data
}

func (p *pipe) GetEndpointChannel(stock string) (chan []value.StockAggregate, error) {
	chans, ok := p.endPoint[stock]
	if !ok {
		return nil, fmt.Errorf("stock do not exist on channel")
	}

	return chans, nil
}