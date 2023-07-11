package polygon

import (
	"log"

	"github.com/Goboolean/fetch-server/internal/infrastructure/ws"
	polygonws "github.com/polygon-io/client-go/websocket"
	"github.com/polygon-io/client-go/websocket/models"
)



func (s *Subscriber) run() {

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
				StockAggsMeta: ws.StockAggsMeta{
					Platform: platformName,
					Symbol:   data.Symbol,
				},

				StockAggsDetail: ws.StockAggsDetail{
					EventType: data.EventType.EventType,
					Min: 		   data.Low,
					Max: 		   data.High,
					Start:     data.Open,
					End:       data.Close,
					StartTime: data.StartTimestamp,
					EndTime:   data.EndTimestamp,
				},
			}

			if err := s.r.OnReceiveStockAggs(stockAggs); err != nil {
				log.Fatal(err)
			}
		}
	}
}


func (s *Subscriber) SubscribeStockAggs(stocks ...string) error {
	return s.conn.Subscribe(polygonws.StocksSecAggs, stocks...)
}

func (s *Subscriber) UnsubscribeStockAggs(stock ...string) error {
	return s.conn.Unsubscribe(polygonws.StocksSecAggs, stock...)
}