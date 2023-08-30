package relay_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Goboolean/fetch-server.v1/cmd/inject"
	"github.com/Goboolean/fetch-server.v1/internal/adapter/websocket"
	"github.com/Goboolean/fetch-server.v1/internal/domain/service/relay"
	"github.com/Goboolean/fetch-server.v1/internal/infrastructure/ws/mock"
)

var instance *relay.Manager

func SetUp() {
	var err error
	ws := websocket.NewMockAdapter().(*websocket.MockAdapter)
	f := mock.New(time.Millisecond*10, ws)

	instance, err = inject.InitMockRelayer(ws)
	if err != nil {
		panic(err)
	}

	if err := ws.RegisterFetcher(f); err != nil {
		panic(err)
	}
	ws.RegisterReceiver(instance)
}

func TearDown() {
	instance.Close()
}

func TestMain(m *testing.M) {
	SetUp()
	code := m.Run()
	TearDown()
	os.Exit(code)
}

func Test_FetchStock(t *testing.T) {

	type args struct {
		stockId string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "FetchStock (case:true)",
			args: args{
				stockId: "stock.google.usa",
			},
			wantErr: false,
		},
		{
			name: "FetchStock (case:false)",
			args: args{
				stockId: "stock.tesla.usa",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := instance.FetchStock(context.Background(), tt.args.stockId); (err != nil) != tt.wantErr {
				t.Errorf("FetchStock() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	var stockId = "stock.google.usa"

	t.Run("IsStockRelayable (case: true)", func(t *testing.T) {
		if relayable := instance.IsStockRelayable(stockId); !relayable {
			t.Errorf("IsStockRelayable() = %v, = %v", relayable, true)
			return
		}
	})

	t.Run("IsStockRelayable (case: false)", func(t *testing.T) {
		if err := instance.StopFetchingStock(context.Background(), stockId); err != nil {
			t.Errorf("StopFetchingStock() = %v", err)
			return
		}

		if relayable := instance.IsStockRelayable(stockId); relayable {
			t.Errorf("IsStockRelayable() = %v, = %v", relayable, false)
			return
		}
	})

}

func Test_Subscribe(t *testing.T) {

	var stockId = "stock.apple.usa"

	var count = 0
	ctx, cancel := context.WithCancel(context.Background())

	t.Run("Subscribe", func(t *testing.T) {

		if err := instance.FetchStock(context.Background(), stockId); err != nil {
			t.Errorf("FetchStock() = %v", err)
			return
		}

		ch, err := instance.Subscribe(ctx, stockId)
		if err != nil {
			t.Errorf("Subscribe() = %v", err)
			return
		}

		go func(ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					return
				case <-ch:
					count++
					break
				}
			}
		}(ctx)

		time.Sleep(time.Millisecond * 100)

		if count <= 5 {
			t.Errorf("Subscribe() = %v, want many", count)
			return
		}
	})

	t.Run("Unsubscribe", func(t *testing.T) {
		cancel()
		var countAfterUnsubscription = count

		time.Sleep(time.Millisecond * 100)

		if countAfterUnsubscription != count {
			t.Errorf("Unsubscribe failed: before: %v, after: %v", countAfterUnsubscription, count)
			return
		}
	})

	t.Run("SubscribeMultiple", func(t *testing.T) {

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		ch1, err := instance.Subscribe(ctx, stockId)
		if err != nil {
			t.Errorf("Subscribe() = %v", err)
			return
		}

		ch2, err := instance.Subscribe(ctx, stockId)
		if err != nil {
			t.Errorf("Subscribe() = %v", err)
			return
		}

		time.Sleep(time.Millisecond * 100)

		cancel()

		if len(ch1) == 0 || len(ch2) == 0 || len(ch1) != len(ch2) || len(ch1) <= 6 || len(ch2) <= 6 {
			t.Errorf("SubscribeMultiple failed: ch1: %v, ch2: %v", len(ch1), len(ch2))
			return
		}
	})

	t.Run("SubscribeBeforeFetching", func(t *testing.T) {

		_, err := instance.Subscribe(context.Background(), "stock.tesla.usa")
		if err != relay.ErrStockNotExists {
			t.Errorf("Subscribe() error = %v, wantErr %v", err, relay.ErrStockNotExists)
			return
		}
	})

	t.Run("StopFetchingAfterSubscribe", func(t *testing.T) {

		ch, err := instance.Subscribe(context.Background(), stockId)
		if err != nil {
			t.Errorf("Subscribe() = %v", err)
			return
		}

		if err := instance.StopFetchingStock(context.Background(), stockId); err != nil {
			t.Errorf("StopFetchingStock() = %v", err)
			return
		}

		time.Sleep(time.Millisecond * 100)

		select {
		case _, ok := <-ch:
			if ok {
				t.Errorf("StopFetchingAfterSubscribe failed: channel is not closed")
				return
			}
		default:
			t.Errorf("StopFetchingAfterSubscribe failed: channel is not closed")
		}
	})
}
