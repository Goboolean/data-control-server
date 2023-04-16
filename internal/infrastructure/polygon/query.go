package polygon

import (
	"github.com/polygon-io/client-go/websocket/models"
	polygonws "github.com/polygon-io/client-go/websocket"
)


func (p *Subscriber) SubscribeStocksSecAggs(stock string) (chan models.EquityAgg, chan error) {
	c := p.conn

	errch := make(chan error)

	if err := c.Subscribe(polygonws.StocksSecAggs, stock); err != nil {
		return nil, errch
	}

	p.ch = make(chan models.EquityAgg, DEFAULT_BUFFER_SIZE)

	go func() {
		defer close(p.ch)

		for {
			select {
			case err := <-c.Error():
				errch <- err
				return
			case out, more := <-c.Output():
				if !more {
					return
				}

				p.ch <- out.(models.EquityAgg)
			}
		}
	}()

	return p.ch, errch
}



