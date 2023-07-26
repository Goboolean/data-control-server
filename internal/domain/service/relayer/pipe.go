package relayer

import (
	"context"
	"reflect"
	"time"

	"github.com/Goboolean/fetch-server/internal/domain/entity"
)

const DEFAULT_BUFFER_SIZE = 10


type pipe struct {
	sinkChan          chan *entity.StockAggregateForm
	filteredChan      chan *entity.StockAggregateForm
	classifiedChanMap map[string]chan *entity.StockAggregate

	connPool map[string] map[int64] conn
}


type conn struct {
	ctx context.Context
	cancel context.CancelFunc
	ch chan *entity.StockAggregate
}

func newConn(ctx context.Context) conn {
	ctx, cancel := context.WithCancel(ctx)
	return conn{
		ctx:    ctx,
		cancel: cancel,
		ch:     make(chan *entity.StockAggregate, DEFAULT_BUFFER_SIZE),
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
				close(conn.ch)
				return
			}
		}
	}(conn.ctx)

	return conn.ch, nil
}



// This method should be executed as goroutine
// It is assured to terminate when channel is closed
func (p *pipe) filterBadTick(in <-chan *entity.StockAggregateForm, out chan<- *entity.StockAggregateForm) {
	for stock := range in {
		if isnil := reflect.DeepEqual(stock, &entity.StockAggregateForm{}); isnil {
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

// This method should be executed as goroutine
// It is assured to terminate when channel is closed
func (p *pipe) relayStockToSubscriber(in <-chan *entity.StockAggregate, out map[int64] conn) {
	for stock := range in {
		for sub := range out {
			out[sub].ch <- stock
		}
	}
}


func newPipe() *pipe {
	return &pipe{
		filteredChan:      make(chan *entity.StockAggregateForm, DEFAULT_BUFFER_SIZE),
		classifiedChanMap: make(map[string]chan *entity.StockAggregate),
		connPool:          make(map[string]map[int64]conn),
		sinkChan:          make(chan *entity.StockAggregateForm, DEFAULT_BUFFER_SIZE),
	}
}


// Run as goroutine, and control lifeccle with ctx.
func (p *pipe) ExecPipe(ctx context.Context) {

	go p.filterBadTick(p.sinkChan, p.filteredChan)
	go p.classifyStock(p.filteredChan, p.classifiedChanMap)
}

func (p *pipe) AddNewPipe(stock string) {
	p.classifiedChanMap[stock] = make(chan *entity.StockAggregate, DEFAULT_BUFFER_SIZE)
	p.connPool[stock] = make(map[int64] conn)
	go p.relayStockToSubscriber(p.classifiedChanMap[stock], p.connPool[stock])
}

func (p *pipe) PlaceOnStartPoint(data *entity.StockAggregateForm) {
	p.sinkChan <- data
}

func (p *pipe) RemovePipe(stock string) {

	for sub := range p.connPool[stock] {
		p.connPool[stock][sub].cancel()
	}
}