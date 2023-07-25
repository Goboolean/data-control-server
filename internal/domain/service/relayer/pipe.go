package relayer

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/Goboolean/fetch-server/internal/domain/entity"
)

const DEFAULT_BUFFER_SIZE = 100


type pipe struct {
	filteredChan      chan *entity.StockAggregateForm
	classifiedChanMap map[string]chan *entity.StockAggregate

	startPoint chan *entity.StockAggregateForm
	endPoint   map[string]chan *entity.StockAggregate

	connPool map[string] map[int64] conn
}


type conn struct {
	ctx context.Context
	ch chan *entity.StockAggregate
}

func newConn(ctx context.Context) conn {
	return conn{
		ctx: ctx,
		ch:  make(chan *entity.StockAggregate, DEFAULT_BUFFER_SIZE),
	}
}


func (p *pipe) getHash() int64 {
	return time.Now().UnixNano()
}


func (p *pipe) RegisterNewSubscriber(ctx context.Context, stockId string) (<-chan *entity.StockAggregate, error) {

	_, ok := p.connPool[stockId]
	if !ok {
		return nil, ErrStockNotExists
	}

	conn := newConn(ctx)
	hash := p.getHash()

	p.connPool[stockId][hash] = conn
	
	go func(ctx context.Context) {
		for {
			select {
			case <- ctx.Done():
				delete(p.connPool[stockId], hash)
				return
			}
		}
	}(ctx)

	return conn.ch, nil
}



// This method should be executed as goroutine
// It is assured to terminate when channel is closed
func (p *pipe) filterBadTick(in <-chan *entity.StockAggregateForm, out chan<- *entity.StockAggregateForm) {
	for stock := range in {
		if isnil := reflect.DeepEqual(stock, &entity.StockAggregateForm{}); isnil {
			fmt.Println(stock.StockID)
			continue
		}
		out <- stock
	}
}

// This method should be executed as goroutine
// It is assured to terminate when channel is closed
func (p *pipe) classifyStock(in <-chan *entity.StockAggregateForm, out map[string]chan *entity.StockAggregate) {
	for stock := range in {
		out[stock.StockID] <- &stock.StockAggregate
	}
}

func (p *pipe) relayStockToSubscriber(in <-chan *entity.StockAggregate, out map[int64] conn) {
	for stock := range in {
		for sub := range out {
			out[sub].ch <- stock
		}
	}
}


func newPipe() *pipe {
	instance := &pipe{
		filteredChan:      make(chan *entity.StockAggregateForm, DEFAULT_BUFFER_SIZE),
		classifiedChanMap: make(map[string]chan *entity.StockAggregate),
		connPool:          make(map[string]map[int64]conn),
		startPoint:        make(chan *entity.StockAggregateForm, DEFAULT_BUFFER_SIZE),
	}

	return instance
}


// Run as goroutine, and control lifeccle with ctx.
func (p *pipe) ExecPipe(ctx context.Context) {

	for stock := range p.endPoint {
		p.classifiedChanMap[stock] = make(chan *entity.StockAggregate, DEFAULT_BUFFER_SIZE)
	}

	go p.filterBadTick(p.startPoint, p.filteredChan)
	go p.classifyStock(p.filteredChan, p.classifiedChanMap)

	go func(ctx context.Context) {
		<-ctx.Done()

		for stock := range p.endPoint {
			close(p.classifiedChanMap[stock])
			delete(p.classifiedChanMap, stock)
		}
		close(p.filteredChan)
	}(ctx)
}

func (p *pipe) AddNewPipe(stock string) {
	p.classifiedChanMap[stock] = make(chan *entity.StockAggregate, DEFAULT_BUFFER_SIZE)
	p.connPool[stock] = make(map[int64] conn)
	go p.relayStockToSubscriber(p.classifiedChanMap[stock], p.connPool[stock])

}

func (p *pipe) PlaceOnStartPoint(data *entity.StockAggregateForm) {
	p.startPoint <- data
}