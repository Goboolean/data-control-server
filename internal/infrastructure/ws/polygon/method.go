package polygon

import (
	"log"

	"github.com/Goboolean/fetch-server/internal/infrastructure/ws"
	polygonws "github.com/polygon-io/client-go/websocket"
	"github.com/polygon-io/client-go/websocket/models"
)



func (s *Subscriber) SubscribeStockAggs(stock string) error {

	if err := s.conn.Subscribe(polygonws.StocksSecAggs, stock); err != nil {
		return err
	}

	return nil
}


func (s *Subscriber) Run() {

	for {
		select {
		case <- s.ctx.Done():
			return
		case err := <-s.conn.Error():
			log.Fatal("error: ", err)
			return
		case out, more := <-s.conn.Output():
			if !more {
				return
			}

			data, ok := out.(models.EquityAgg)
			if !ok {
				log.Fatal("failed to cast data")
				return
			}

			stockAggs := &ws.StockAggregate{
				EventType: data.EventType.EventType,
				Min: 		   data.Low,
				Max: 		   data.High,
				Start:     data.Open,
				End:       data.Close,
				StartTime: data.StartTimestamp,
				EndTime:   data.EndTimestamp,
			}

			if err := s.r.OnReceiveStockAggs(stockAggs); err != nil {
				log.Fatal(err)
			}
		}
	}
}



func (s *Subscriber) UnsubscribeStockAggs(stock string) error {
	return nil
}