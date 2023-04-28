package server

import (
	"fmt"
	"net"
	"os"

	adapter "github.com/Goboolean/stock-fetch-server/internal/adapter/grpc"
	"github.com/Goboolean/stock-fetch-server/internal/infrastructure/grpc/config"
	"google.golang.org/grpc"
)

var (
	FETCH_SERVER_PORT = os.Getenv("FETCH_SERVER_PORT")
)

func Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", FETCH_SERVER_PORT))

	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	config.RegisterStockConfiguratorServer(grpcServer, adapter.New())

	if err := grpcServer.Serve(lis); err != nil {
		return err
	}

	return nil
}
