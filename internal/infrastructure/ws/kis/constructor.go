package kis

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/Goboolean/fetch-server/internal/infrastructure/ws"
	"github.com/Goboolean/shared/pkg/resolver"
	"github.com/gorilla/websocket"
)

const platformName = "kis"

const address = "ops.koreainvestment.com:21000"

type Data struct {
	GrantType string `json:"grant_type"`
	AppKey    string `json:"appkey"`
	SecretKey string `json:"secretkey"`
}

type Response struct {
	ApprovalKey string `json:"approval_key"`
}

type Subscriber struct {
	conn *websocket.Conn

	approval_key string

	ctx    context.Context
	cancel context.CancelFunc
	r      ws.Receiver
}

func (s *Subscriber) getApprovalKey(Appkey string, Secretkey string) (string, error) {
	data := &Data{
		GrantType: "client_credentials",
		AppKey:    Appkey,
		SecretKey: Secretkey,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	response, err := http.Post("https://openapi.koreainvestment.com:9443/oauth2/Approval", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var res *Response

	if err := json.Unmarshal(body, &res); err != nil {
		return "", err
	}

	return res.ApprovalKey, nil
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

	appkey, err := c.GetStringKey("KIS_APPKEY")
	if err != nil {
		panic(err)
	}

	secretkey, err := c.GetStringKey("KIS_SECRET")
	if err != nil {
		panic(err)
	}

	approval_key, err := instance.getApprovalKey(appkey, secretkey)
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
	return handler("")
}
