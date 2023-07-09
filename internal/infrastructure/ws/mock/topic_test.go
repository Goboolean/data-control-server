package mock

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/Goboolean/fetch-server/internal/infrastructure/ws"
)


var (
	generater *mockGenerater

	topic = "test"
	ch chan *ws.StockAggregate
)

func SetupMockGenerater() {
	ctx := context.Background()
	ch = make(chan *ws.StockAggregate)
	
	duration := time.Second / 10 // the data is generated every 0.1 second in average.
	
	generater = newMockGenerater(topic, ctx, ch, duration)
}

func TeardownMockGenerater() {
	close(ch)
	generater.Close()
}



// It fails the test when it generates empty data.
func Test_generateRandomStockAggs(t *testing.T) {
	
	SetupMockGenerater()

	agg := generater.generateRandomStockAggs()
	if same := reflect.DeepEqual(agg, ws.StockAggregate{}); same {
		t.Errorf("generateRandomStockAggs() = %v, want not empty", agg)
	}

	TeardownMockGenerater()
}



// It verdicts the test as success when it generates data 5 times for a second.
func Test_newMockGenerater(t *testing.T) {
	
	SetupMockGenerater()

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second))
	defer cancel()

	for count := 5; count >= 0; count-- {
		select {
		case <- ctx.Done():
			t.Errorf("newMockGenerater() got timeout")
		case <- ch:
			continue
		}
	}

	generater.Close()

	// empty channel to test that it does not generate data after closing.
	for len(ch) > 0 {
		<- ch
	}

	select {
	case <- ch:
		t.Errorf("newMockGenerater got data after closing")
	case <- time.After(time.Second / 10):
		break
	}

	TeardownMockGenerater()
}
