package grpc_test

import (
	"context"
	"os"
	"testing"
	"time"

	pb "github.com/Goboolean/fetch-server.v1/api/grpc"
	grpc_adapter "github.com/Goboolean/fetch-server.v1/internal/adapter/grpc"
	server "github.com/Goboolean/fetch-server.v1/internal/infrastructure/grpc"
	"github.com/Goboolean/shared/pkg/resolver"

	_ "github.com/Goboolean/fetch-server.v1/internal/util/env"
)

var (
	instance *server.Host
	client   *server.Client
)

func SetUp() {
	var err error

	instance, err = server.New(&resolver.ConfigMap{
		"PORT": os.Getenv("SERVER_PORT"),
	}, grpc_adapter.NewMockAdapter())
	if err != nil {
		panic(err)
	}

	time.Sleep(1 * time.Second)

	client, err = server.NewClient(&resolver.ConfigMap{
		"PORT": os.Getenv("SERVER_PORT"),
	})
	if err != nil {
		panic(err)
	}
}

func TearDown() {
	instance.Close()
	client.Close()
}

func TestMain(m *testing.M) {
	SetUp()
	code := m.Run()
	TearDown()
	os.Exit(code)
}

func Test_Constructur(t *testing.T) {

	t.Run("UpdateStockConfigOne", func(t *testing.T) {
		_, err := client.GetStockConfigAll(context.Background(), &pb.Null{})
		if err != nil {
			t.Errorf("Ping() error = %v", err)
			return
		}
	})

	t.Run("UpdateStockConfigMany", func(t *testing.T) {
		_, err := client.GetStockConfigAll(context.Background(), &pb.Null{})
		if err != nil {
			t.Errorf("Ping() error = %v", err)
			return
		}
	})

	t.Run("GetStockConfigOne", func(t *testing.T) {
		_, err := client.GetStockConfigAll(context.Background(), &pb.Null{})
		if err != nil {
			t.Errorf("Ping() error = %v", err)
			return
		}
	})

	t.Run("GetStockConfigAll", func(t *testing.T) {
		_, err := client.GetStockConfigAll(context.Background(), &pb.Null{})
		if err != nil {
			t.Errorf("Ping() error = %v", err)
			return
		}
	})
}
