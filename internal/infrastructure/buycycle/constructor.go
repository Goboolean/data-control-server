package buycycle

import (
	"context"
	"fmt"
	"net/url"

	"github.com/Goboolean/fetch-server/internal/infrastructure/ws"
	"github.com/Goboolean/shared/pkg/resolver"
	"github.com/gorilla/websocket"
)

var DEFAULT_BUFFER_SIZE = 1000

type Subscriber struct {
	*websocket.Conn

	ctx    context.Context
	cancel context.CancelFunc
	r      ws.Receiver
}

func New(c *resolver.ConfigMap, r ws.Receiver) *Subscriber {

	host, err := c.GetStringKey("HOST")
	if err != nil {
		panic(err)
	}

	port, err := c.GetStringKey("PORT")
	if err != nil {
		panic(err)
	}

	path, err := c.GetStringKey("PATH")
	if err != nil {
		panic(err)
	}

	address := fmt.Sprintf("%s:%s", host, port)

	u := url.URL{
		Scheme: "ws",
		Host:   address,
		Path:   path,
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	instance := &Subscriber{
		Conn:   conn,
		ctx:    ctx,
		cancel: cancel,
		r:      r,
	}

	go RelayMessageToReceiver(instance)

	return instance
}


func (s *Subscriber) Close() error {
	
	if err := s.Conn.Close(); err != nil {
		return err
	}

	s.cancel()
	return nil
}


func (s *Subscriber) Ping() error {
	handler := s.Conn.PingHandler()
	return handler("")
}