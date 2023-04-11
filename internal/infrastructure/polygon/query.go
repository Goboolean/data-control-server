package polygon

import (
	"github.com/polygon-io/client-go/websocket/models"
	polygonws "github.com/polygon-io/client-go/websocket"
)


func (p *PolygonWs) SubscribeStocksSecAggs(stock string) (chan models.EquityAgg, error) {
	c := p.conn

	if err := c.Subscribe(polygonws.StocksSecAggs, stock); err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			case err := <-c.Error():
				panic(err)
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



