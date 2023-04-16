package buycycle

import "encoding/json"

func (w *Subscriber) SubscribeStocks(stock string) (chan StockAggregate, chan error) {

	w.ch = make(chan StockAggregate, DEFAULT_BUFFER_SIZE)

	errCh := make(chan error)

	go func() {
		defer close(w.ch)

		for {
			_, message, err := w.conn.ReadMessage()

			if err != nil {
				errCh <- err
				return
			}

			var data StockAggregate

			if err := json.Unmarshal(message, &data); err != nil {
				errCh <- err
				return
			}

			w.ch <- data
		}
	}()

	return w.ch, errCh
}
