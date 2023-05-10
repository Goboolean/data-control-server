package polygon

import (
	"context"

	"github.com/Goboolean/shared-packages/pkg/resolver"
	polygonws "github.com/polygon-io/client-go/websocket"
)


var DEFAULT_BUFFER_SIZE = 1000


type Subscriber struct {
	*polygonws.Client

	r Receiver
	ctx context.Context
	cancel context.CancelFunc
}



func New(c *resolver.Config, r Receiver) *Subscriber {

	if err := c.ShouldPWExist(); err != nil {
		panic(err)
	}

	conn, err := polygonws.New(polygonws.Config{
		APIKey:    c.Password,
		Feed:      polygonws.RealTime,
		Market:    polygonws.Stocks,
	})
	
	if err != nil {
		panic(err)
	}

	if err := conn.Connect(); err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	instance := &Subscriber{
		Client: conn,
		r: r,
		ctx: ctx,
		cancel: cancel,
	}

	go RelayMessageToReceiver(instance)

	return instance
}



func (s *Subscriber) Close() error {
	if err := s.Close(); err != nil {
		return err
	}

	s.cancel()

	return nil
}