package buycycle

import (
	"context"
	"fmt"
	"net/url"

	"github.com/Goboolean/shared-packages/pkg/resolver"
	"github.com/gorilla/websocket"
)

var DEFAULT_BUFFER_SIZE = 1000

type Subscriber struct {
	*websocket.Conn

	ctx context.Context
	cancel context.CancelFunc
	r Receiver
}

func New(c *resolver.Config, r Receiver) *Subscriber {

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

	instance := &Subscriber{
		Conn: conn,
		ctx: ctx,
		cancel: cancel,
		r: r,
	}

	go RelayMessageToReceiver(instance);

	return instance
}


func (s *Subscriber) Close() error {
	
	if err := s.Close(); err != nil {
		return err
	}

	s.cancel()
	return nil
}
