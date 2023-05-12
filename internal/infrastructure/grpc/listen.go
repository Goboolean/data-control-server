package server

import (
	"context"
	"fmt"
	"net"
	"os"

	api "github.com/Goboolean/stock-fetch-server/api/grpc"
	adapter "github.com/Goboolean/stock-fetch-server/internal/adapter/grpc"
	"google.golang.org/grpc"
)

var (
	FETCH_SERVER_PORT = os.Getenv("FETCH_SERVER_PORT")
)

func Run(ctx context.Context, ch chan error, adapter *adapter.StockConfiguratorAdapter) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", FETCH_SERVER_PORT))

	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	api.RegisterStockConfiguratorServer(grpcServer, adapter)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			ch <- err
		}
	}()

	select {
	case <- ctx.Done():
		grpcServer.Stop()
	}
}
