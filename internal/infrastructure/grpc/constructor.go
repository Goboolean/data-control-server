package grpc

import (
	"fmt"
	"net"
	"sync"

	api "github.com/Goboolean/fetch-server/api/grpc"
	"github.com/Goboolean/shared/pkg/resolver"
	"google.golang.org/grpc"
)

type Host struct {
	server *grpc.Server
}

var (
	instance *Host
	once     sync.Once
)

func New(c *resolver.ConfigMap, adapter api.StockConfiguratorServer) *Host {

	once.Do(func() {
		port, err :=  c.GetStringKey("PORT")
		if err != nil {
			panic("fetch server port required")
		}
	
		lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
		if err != nil {
			panic(err)
		}
	
		grpcServer := grpc.NewServer()
		api.RegisterStockConfiguratorServer(grpcServer, adapter)
		go grpcServer.Serve(lis)

		instance = &Host{
			server: grpcServer,
		}
	})

	return instance
}



func (s *Host) Close() {
	s.server.GracefulStop()
}
