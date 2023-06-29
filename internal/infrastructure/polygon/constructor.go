package polygon

import (
	"context"
	"sync"
	"time"

	"github.com/Goboolean/shared-packages/pkg/resolver"
	polygonws "github.com/polygon-io/client-go/websocket"
)


var DEFAULT_BUFFER_SIZE = 1000


type Subscriber struct {
	conn *polygonws.Client

	r Receiver
	ctx context.Context
}

var (
	instance *Subscriber
	once sync.Once
)



func New(c *resolver.Config, r Receiver) *Subscriber {

	if err := c.ShouldPWExist(); err != nil {
		panic(err)
	}

	once.Do(func() {
		conn, err := polygonws.New(polygonws.Config{
			APIKey:    c.Password,
			Feed:      polygonws.RealTime,
			Market:    polygonws.Stocks,
		})
		
		if err != nil {
			panic(err)
		}
	
		instance = &Subscriber{
			conn: conn,
			r: r,
		}
	
	})

	return instance
}


func (p *Subscriber) tryRun() error {
	if err := p.conn.Connect(); err != nil {
		return err
	}

	go RelayMessageToReceiver(p)
	return nil
}


func (s *Subscriber) Run() {
	go func() {
		for {
			if err := s.tryRun(); err != nil {
				time.Sleep(time.Hour)
			} else {
				break
			}
		}
	}()
}

func (s *Subscriber) Close() error {
	if err := s.Close(); err != nil {
		return err
	}

	return nil
}