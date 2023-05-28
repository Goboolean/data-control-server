package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"sync"

	api "github.com/Goboolean/stock-fetch-server/api/grpc"
	adapter "github.com/Goboolean/stock-fetch-server/internal/adapter/grpc"
	"google.golang.org/grpc"
)



type Host struct {
	server *grpc.Server
	impl *adapter.StockConfiguratorAdapter
}


var (
	instance *Host
	once sync.Once
)

func New(adapter *adapter.StockConfiguratorAdapter) *Host {
	once.Do(func() {
		instance = &Host{
			impl: adapter,
		}		
	})

	return instance
}




func (h *Host) Run(ctx context.Context, adapter *adapter.StockConfiguratorAdapter) {

	port, flag := os.LookupEnv("FETCH_SERVER_PORT")

	if !flag {
		panic("fetch server port required")
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))

	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	api.RegisterStockConfiguratorServer(grpcServer, adapter)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			panic(err)
		}
	}()
}



func (s *Host) Close() {
	s.server.GracefulStop()
}