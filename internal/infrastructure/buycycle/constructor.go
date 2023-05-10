package buycycle

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/Goboolean/shared-packages/pkg/resolver"
	"github.com/gorilla/websocket"
)

var DEFAULT_BUFFER_SIZE = 1000

type Client struct {
	*websocket.Conn

	ctx context.Context
	cancel context.CancelFunc
}

func New(c *resolver.Config, r Receiver) *Client {

	if err := c.ShouldHostExist(); err != nil {
		panic(err)
	}

	if err := c.ShouldPortExist(); err != nil {
		panic(err)
	}

	c.Address = fmt.Sprintf("%s:%s", c.Host, c.Port)

	u := url.URL{
		Scheme: "ws",
		Host:   c.Address,
		Path:   "BUYCYCLE_PATH",
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func(ctx context.Context) {

		var data StockAggregate

		for {
			select {
			case <- ctx.Done():
				return
			default:

				if err := conn.ReadJSON(&data); err != nil {
					log.Fatalf("failed to read json data: %v", err)
					continue
				}

				if err := r.OnReceiveBuycycleStockAggs(data); err != nil {
					log.Fatalf("failed to process data: %v", err)
					continue
				}
			}
		}
	}(ctx)

	return &Client{
		Conn: conn,
		ctx: ctx,
		cancel: cancel,
	}
}


func (s *Client) Close() error {
	
	if err := s.Close(); err != nil {
		return err
	}

	s.cancel()
	return nil
}
