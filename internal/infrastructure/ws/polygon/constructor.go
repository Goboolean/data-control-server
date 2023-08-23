package polygon

import (
	"context"

	"github.com/Goboolean/fetch-server/internal/infrastructure/ws"
	"github.com/Goboolean/shared/pkg/resolver"
	polygonws "github.com/polygon-io/client-go/websocket"
)

const platformName = "polygon"

type Subscriber struct {
	conn *polygonws.Client

	r      ws.Receiver
	ctx    context.Context
	cancel context.CancelFunc
}

func New(c *resolver.ConfigMap, r ws.Receiver) (*Subscriber, error) {

	key, err := c.GetStringKey("KEY")
	if err != nil {
		return nil, err
	}

	conn, err := polygonws.New(polygonws.Config{
		APIKey: key,
		Feed:   polygonws.RealTime,
		Market: polygonws.Stocks,
	})

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	instance := &Subscriber{
		conn:   conn,
		r:      r,
		ctx:    ctx,
		cancel: cancel,
	}

	go instance.run()
	return instance, nil
}

func (s *Subscriber) PlatformName() string {
	return platformName
}

func (s *Subscriber) Close() error {
	s.cancel()
	s.conn.Close()
	return nil
}

func (s *Subscriber) Ping() error {
	// Ping() does not use directly *polygon.Client.
	// TODO: Find a way to check connection is alive.

	return nil
}
