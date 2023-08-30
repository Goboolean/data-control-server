package grpc

import (
	"fmt"
	"net"

	api "github.com/Goboolean/fetch-server.v1/api/grpc"
	"github.com/Goboolean/shared/pkg/resolver"
	"google.golang.org/grpc"
)

type Host struct {
	server *grpc.Server
}

func New(c *resolver.ConfigMap, adapter api.StockConfiguratorServer) (*Host, error) {

	port, err := c.GetStringKey("PORT")
	if err != nil {
		return nil, err
	}

	address := fmt.Sprintf(":%s", port)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	grpcServer := grpc.NewServer()
	api.RegisterStockConfiguratorServer(grpcServer, adapter)
	go grpcServer.Serve(lis)

	return &Host{
		server: grpcServer,
	}, nil
}

func (s *Host) Close() {
	s.server.GracefulStop()
}
