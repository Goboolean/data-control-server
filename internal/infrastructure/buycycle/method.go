package buycycle

import "log"






func (s *Subscriber) SubscribeStockAggs(stock string) error {
	return s.WriteJSON(stock);
}


func RelayMessageToReceiver(s *Subscriber) {
	var data StockAggregate

	for {
		select {
		case <- s.ctx.Done():
			return
		default:

			if err := s.ReadJSON(&data); err != nil {
				log.Fatalf("failed to read json data: %v", err)
				continue
			}

			if err := s.r.OnReceiveBuycycleStockAggs(data); err != nil {
				log.Fatalf("failed to process data: %v", err)
				continue
			}
		}
	}
}