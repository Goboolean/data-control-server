package buycycle

import (
	"context"
	"fmt"
	"net/url"

	"github.com/Goboolean/fetch-server/internal/infrastructure/ws"
	"github.com/Goboolean/shared/pkg/resolver"
	"github.com/gorilla/websocket"
)


const platformName = "buycycle"



type Subscriber struct {
	conn *websocket.Conn

	ctx    context.Context
	cancel context.CancelFunc
	r      ws.Receiver
}



func New(c *resolver.ConfigMap, ctx context.Context, r ws.Receiver) *Subscriber {

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

	ctx, cancel := context.WithCancel(ctx)

	instance := &Subscriber{
		conn:   conn,
		ctx:    ctx,
		cancel: cancel,
		r:      r,
	}

	go instance.run()

	return instance
}


func (s *Subscriber) PlatformName() string {
	return platformName
}


func (s *Subscriber) Close() error {
	
	if err := s.conn.Close(); err != nil {
		return err
	}

	s.cancel()
	return nil
}


func (s *Subscriber) Ping() error {
	handler := s.conn.PingHandler()
	return handler("")
}