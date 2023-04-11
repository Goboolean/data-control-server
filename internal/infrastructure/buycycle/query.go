package buycycle

import "encoding/json"



func (w *BuycycleWs) SubscribeStocks(stock string) (chan StockAggregate, chan error) {

	w.ch = make(chan StockAggregate, DEFAULT_BUFFER_SIZE)

	errch := make(chan error)

	go func() {
		defer close(w.ch)

		for {
			_, message, err := w.conn.ReadMessage()

			if err != nil {
				errch <- err
				return
			}

			var data StockAggregate

			if err := json.Unmarshal(message, &data); err != nil {
				errch <- err
				return
			}

			w.ch <- data
		}
	}()

	return w.ch, errch
}