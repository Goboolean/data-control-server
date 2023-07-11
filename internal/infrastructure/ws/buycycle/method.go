package buycycle

import (
	"log"

	"github.com/Goboolean/fetch-server/internal/infrastructure/ws"
)

// It is not sure whether rolling back subscription is provided on buycycle opensource.
// Therefore, this package is not implemented yet.

func (s *Subscriber) run() {
	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			
		}

		var agg *ws.StockAggregate
		if err := s.conn.ReadJSON(agg); err != nil {
			log.Fatal(err)
		}

		if err := s.r.OnReceiveStockAggs(agg); err != nil {
			log.Fatal(err)
		}
	}
}



func (s *Subscriber) SubscribeStockAggs(symbols ...string) error {
	for _, symbol := range symbols {
		req := &RequestJson{
			Header: HeaderJson{},
			Body: RequestBodyJson{	
				Query: struct {Shcode string `json:"shcode"`} {
					Shcode: symbol,
				},
				// TODO: Add more fields
			},
		}

		if err := s.conn.WriteJSON(req); err != nil {
			return err
		}
	}

	return nil
}



func (s *Subscriber) UnsubscribeStockAggs(stock ...string) error {
	// TODO: check how Unsubscribe works on buycycle than implement it.
	return nil
}