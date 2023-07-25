package relayer_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Goboolean/fetch-server/internal/adapter/meta"
	"github.com/Goboolean/fetch-server/internal/adapter/persistence"
	"github.com/Goboolean/fetch-server/internal/adapter/transaction"
	"github.com/Goboolean/fetch-server/internal/adapter/websocket"
	"github.com/Goboolean/fetch-server/internal/domain/service/relayer"
	"github.com/Goboolean/fetch-server/internal/infrastructure/ws/mock"
)

var instance *relayer.RelayerManager



func NewMock() *relayer.RelayerManager {
	
	var (
		db           = persistence.NewMockAdapter()
		tx           = transaction.NewMock()
		meta         = meta.NewMockAdapter()
		ws = websocket.NewAdapter()
		f = mock.New(context.Background(), time.Millisecond * 10, ws)
	)

	if err := ws.RegisterFetcher(f); err != nil {
		panic(err)
	}

	instance = relayer.New(db, tx, meta, ws)
	ws.RegisterReceiver(instance)

	return instance
}



func SetUp() {
	instance = NewMock()
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
			name: "existing stock",
			args: args{
				stockId: "stock.google.usa",
			},
			wantErr: false,
		},
		{
			name: "non existing stock",
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

	if relayable := instance.IsStockRelayable(stockId); !relayable {
		t.Errorf("IsStockRelayable() = %v, = %v", relayable, true)
		return
	}

	if err := instance.StopFetchingStock(context.Background(), stockId); err != nil {
		t.Errorf("StopFetchingStock() = %v", err)
		return
	}

	if relayable := instance.IsStockRelayable(stockId); relayable {
		t.Errorf("IsStockRelayable() = %v, = %v", relayable, false)
		return
	}
}


func Test_Subscribe(t *testing.T) {

	var stockId = "stock.apple.usa"

	if err := instance.FetchStock(context.Background(), stockId); err != nil {
		t.Errorf("FetchStock() = %v", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch, err := instance.Subscribe(ctx, stockId)
	if err != nil {
		t.Errorf("Subscribe() = %v", err)
		return
	}

	var count = 0

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

	if (count <= 5) {
		t.Errorf("Subscribe() = %v, want many", count)
		return
	}
}
