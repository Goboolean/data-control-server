package polygon

import (
	"log"

	polygonws "github.com/polygon-io/client-go/websocket"
	"github.com/polygon-io/client-go/websocket/models"
)


func (p *Subscriber) SubscribeStocksSecAggs(stock string) (<-chan models.EquityAgg, error) {
	c := p.conn

	if err := c.Subscribe(polygonws.StocksSecAggs, stock); err != nil {
		return nil, err
	}

	p.ch = make(chan models.EquityAgg, DEFAULT_BUFFER_SIZE)

	go func() {
		defer close(p.ch)

		for {
			select {
			case err := <-c.Error():
				log.Fatal("error: ", err)
				return
			case out, more := <-c.Output():
				if !more {
					return
				}

				p.ch <- out.(models.EquityAgg)
			}
		}
	}()

	return p.ch, nil
}



