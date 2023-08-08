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

	approvalKey string

	ctx    context.Context
	cancel context.CancelFunc
	r      ws.Receiver

	subscribed chan string
}



func New(c *resolver.ConfigMap, r ws.Receiver) *Subscriber {

	u := url.URL{Scheme: "ws", Host: address}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	instance := &Subscriber{
		conn:   conn,
		ctx:    ctx,
		cancel: cancel,
		r:      r,
		subscribed: make(chan string),
	}

	appkey, err := c.GetStringKey("APPKEY")
	if err != nil {
		panic(err)
	}

	secretkey, err := c.GetStringKey("SECRET")
	if err != nil {
		panic(err)
	}

	approvalKey, err := instance.GetApprovalKey(appkey, secretkey)
	if err != nil {
		panic(err)
	}

	instance.approvalKey = approvalKey

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
	close(s.subscribed)
	return nil
}

func (s *Subscriber) Ping() error {
	handler := s.conn.PingHandler()
	return handler("Ping")
}
