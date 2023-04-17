package buycycle

import (
	"fmt"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
)

var (
	BUYCYCLE_HOST = os.Getenv("BUYCYCLE_HOST")
	BUYCYCLE_PORT = os.Getenv("BUYCYCLE_PORT")
	BUYCYCLE_PATH = os.Getenv("BUYCYCLE_PATH")

	DEFAULT_BUFFER_SIZE = 1000
)

type Subscriber struct {
	conn *websocket.Conn
	ch   chan StockAggregate
}

func New() *Subscriber {

	u := url.URL{
		Scheme: "ws",
		Host:   fmt.Sprintf("%s:%s", BUYCYCLE_HOST, BUYCYCLE_PORT),
		Path:   BUYCYCLE_PATH,
	}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	if err != nil {
		panic(err)
	}

	return &Subscriber{
		conn: c,
		ch:   make(chan StockAggregate, DEFAULT_BUFFER_SIZE),
	}
}
