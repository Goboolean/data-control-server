package relayer

import (
	"context"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/Goboolean/fetch-server/internal/domain/entity"
)



func generateRandomStockAggregate() entity.StockAggregate {

	rand.Seed(time.Now().UnixNano())

	return entity.StockAggregate{
		EventType: "stock",
		Average:  1.0 + rand.Float64()*(2.0),
		Min:      1.0 + rand.Float64()*(2.0),
		Max:      1.0 + rand.Float64()*(2.0),
		Start: 	  1.0 + rand.Float64()*(2.0),
		End:      1.0 + rand.Float64()*(2.0),
		StartTime: time.Now().Unix(),
		EndTime:  time.Now().Unix(),		
	}
}


func generateRandomStockAggregateForm(stockId string) entity.StockAggregateForm {

	agg := generateRandomStockAggregate()

	return entity.StockAggregateForm{
		StockAggsMeta: entity.StockAggsMeta{
			StockID: stockId,
		},
		StockAggregate: agg,
	}
}

func Test_generateRandomStockAggregateForm(t *testing.T) {
	
	var stockId = "test"
	var count = 2

	agg := generateRandomStockAggregateForm(stockId)

	if equals := reflect.DeepEqual(agg, &entity.StockAggregateForm{}); equals {
		t.Error("generateRandomStockAggregateForm() failed: got nil stockAggs")
		return
	}

	ch := make(chan *entity.StockAggregateForm, DEFAULT_BUFFER_SIZE)

	for i := 0; i < count; i++ {
		agg := generateRandomStockAggregateForm(stockId)
		ch <- &agg
	}

	received1 := <- ch
	received2 := <- ch

	if equals := reflect.DeepEqual(received1, received2); equals {
		t.Error("generateRandomStockAggregateForm() failed: got same stockAggs")
		return
	}
}



func Test_filterBadTick(t *testing.T) {

	var count = 5
	var stockId = "test"

	
	var p *pipe

	sink := make(chan *entity.StockAggregateForm, DEFAULT_BUFFER_SIZE)
	source := make(chan *entity.StockAggregateForm, DEFAULT_BUFFER_SIZE)
	defer close(sink)
	defer close(source)

	go p.filterBadTick(sink, source)

	// generate vague data

	sink <- &entity.StockAggregateForm{}

	time.Sleep(time.Millisecond * 10)

	select {
	case <- source:
		t.Error("filterBadTick() failed: failed to filter out vague data (bad tick)")
		return
	default:
		break
	}

	// generate existable data

	for i := 0; i < count; i++ {
		agg := generateRandomStockAggregateForm(stockId)
		sink <- &agg
	}

	time.Sleep(time.Millisecond * 10)

	for i := 0; i < count; i++ {
		select {
		case <- source:
			break
		default:
			t.Error("filterBadTick() failed: failed to transmit data from startPoint to endpoint")
			return
		}
	}

	select {
	case <- source:
		t.Error("filterBadTick() failed: endPoint received non existent data")
		return
	default:
		break
	}

}


func Test_classifyStock(t *testing.T) {

	var count = 5
	var targetStockId = "test"
	var antiTargetStockId = "test2"

	var p *pipe

	sink := make(chan *entity.StockAggregateForm, DEFAULT_BUFFER_SIZE)
	source := make(map[string] chan *entity.StockAggregate)
	source[targetStockId] = make(chan *entity.StockAggregate, DEFAULT_BUFFER_SIZE)
	source[antiTargetStockId] = make(chan *entity.StockAggregate, DEFAULT_BUFFER_SIZE)

	defer func() {
		close(sink)
		for k, v := range source {
			close(v)
			delete(source, k)
		}
	}()

	go p.classifyStock(sink, source)

	for i := 0; i < count; i++ {
		agg := generateRandomStockAggregateForm(targetStockId)
		sink <- &agg
	}

	time.Sleep(time.Millisecond)

	select {
		case <- source[antiTargetStockId]:
			t.Error("received non existent data")
			return
		default:
			break
	}

	for i := 0; i < count; i++ {
		select {
		case <- source[targetStockId]:
			break
		default:
			t.Error("filterBadTick() failed: failed to transmit data from startPoint to endpoint")
			return
		}
	}

	select {
	case <- source[targetStockId]:
		t.Error("filterBadTick() failed: endPoint received non existent data")
		return
	default:
		break
	}

}


func Test_relayStockToSubscriber(t *testing.T) {

	var count = 5

	var p *pipe

	sink := make(chan *entity.StockAggregate, DEFAULT_BUFFER_SIZE)
	source := make(map[int64]conn)

	source[0] = newConn(context.Background())

	go p.relayStockToSubscriber(sink, source)

	source[1] = newConn(context.Background())

	ch1 := source[0].ch
	ch2 := source[1].ch

	select {
	case <- ch1:
		t.Error("filterBadTick() failed: source received non existent data")
		return
	default:
		break
	}

	for i := 0; i < count; i++ {
		agg := generateRandomStockAggregate()
		sink <- &agg
	}

	time.Sleep(10 * time.Millisecond)

	for i := 0; i < count; i++ {
		select {
		case <- ch1:
			break
		default:
			t.Error("filterBadTick() failed: failed to transmit data from sink to source")
			return
		}
	}

	select {
	case <- ch1:
		t.Error("filterBadTick() failed: source received non existent data")
		return
	default:
		break
	}

	for i := 0; i < count; i++ {
		select {
		case <- ch2:
			break
		default:
			t.Error("filterBadTick() failed: failed to transmit data from sink to source")
			return
		}
	}

}


func Test_pipe(t *testing.T) {

	p := newPipe()

	var stockId = "test"
	var count = 5

	p.ExecPipe(context.Background())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	p.AddNewPipe(stockId)

	ch, err := p.RegisterNewSubscriber(ctx, stockId)
	if err != nil {
		t.Errorf("RegisterNewSubscriber() failed: %v", err)
		return
	}

	for i := 0; i < count; i++ {
		agg := generateRandomStockAggregateForm(stockId)
		p.PlaceOnStartPoint(&agg)
	}

	time.Sleep(10 * time.Millisecond)


	for i := 0; i < count; i++ {
		select {
		case <- ch:
			break
		default:
			t.Error("filterBadTick() failed: failed to transmit data from sink to source")
			return
		}
	}

	cancel()

	agg := generateRandomStockAggregateForm(stockId)
	p.PlaceOnStartPoint(&agg)

	time.Sleep(10 * time.Millisecond)


	select {
	case 	_, ok := <- ch:
		if ok {
			t.Error("RegisterNewSubscriber() failed: channel should be closed after cancel() is called")
			return
		}
		break
	default:
		break
	}

}