package grpc_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	pb "github.com/Goboolean/fetch-server/api/grpc"
	server "github.com/Goboolean/fetch-server/internal/infrastructure/grpc"
	grpc_adapter "github.com/Goboolean/fetch-server/internal/adapter/grpc"
	"github.com/Goboolean/fetch-server/internal/util/env"
	"github.com/Goboolean/shared/pkg/resolver"
	"github.com/joho/godotenv"

	"google.golang.org/grpc"
)




var (
	instance *server.Host
	client pb.StockConfiguratorClient
)



func NewClient() pb.StockConfiguratorClient {

	address := fmt.Sprintf(":%s", os.Getenv("SERVER_PORT"))
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		panic(err)
	}

	return pb.NewStockConfiguratorClient(conn)
}


func SetUp() {

	if err := os.Chdir(env.Root); err != nil {
		panic(err)
	}

	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	instance = server.New(&resolver.ConfigMap{
		"PORT": os.Getenv("SERVER_PORT"),
	}, grpc_adapter.NewMockAdapter())

	time.Sleep(1 * time.Second)

	client = NewClient()
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