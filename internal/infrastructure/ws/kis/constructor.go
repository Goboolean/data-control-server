package kis

import (
	"context"

	"net/url"

	"github.com/Goboolean/fetch-server/internal/infrastructure/ws"
	"github.com/Goboolean/shared/pkg/resolver"
	"github.com/gorilla/websocket"
)

const platformName = "kis"

const address = "ops.koreainvestment.com:21000"


type Subscriber struct {
	conn *websocket.Conn

	approval_key string

	ctx    context.Context
	cancel context.CancelFunc
	r      ws.Receiver
}



func New(c *resolver.ConfigMap, ctx context.Context, r ws.Receiver) *Subscriber {

	u := url.URL{Scheme: "ws", Host: address}

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

	appkey, err := c.GetStringKey("APPKEY")
	if err != nil {
		panic(err)
	}

	secretkey, err := c.GetStringKey("SECRET")
	if err != nil {
		panic(err)
	}

	approval_key, err := instance.GetApprovalKey(appkey, secretkey)
	if err != nil {
		panic(err)
	}

	instance.approval_key = approval_key

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
	return handler("Ping")
}
