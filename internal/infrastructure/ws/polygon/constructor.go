package polygon

import (
	"context"
	"sync"

	"github.com/Goboolean/fetch-server/internal/infrastructure/ws"
	"github.com/Goboolean/shared/pkg/resolver"
	polygonws "github.com/polygon-io/client-go/websocket"
)

var DEFAULT_BUFFER_SIZE = 1000

type Subscriber struct {
	conn *polygonws.Client

	r   ws.Receiver
	ctx context.Context
	cancel context.CancelFunc
}

var (
	instance *Subscriber
	once     sync.Once
)

func New(c *resolver.ConfigMap, ctx context.Context, r ws.Receiver) *Subscriber {

	key, err := c.GetStringKey("KEY")
	if err != nil {
		panic(err)
	}

	once.Do(func() {
		conn, err := polygonws.New(polygonws.Config{
			APIKey: key,
			Feed:   polygonws.RealTime,
			Market: polygonws.Stocks,
		})

		if err != nil {
			panic(err)
		}

		ctx, cancel := context.WithCancel(ctx)

		instance = &Subscriber{
			conn: conn,
			r:    r,
			ctx:  ctx,
			cancel: cancel,
		}
	})

	return instance
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