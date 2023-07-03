package transmission

import (
	"context"
	"log"

	"github.com/Goboolean/fetch-server/internal/adapter/transaction"
)

func (t *Transmitter) SubscribeRelayer(stock string) error {
	ch, err := t.relayer.Subscribe(stock)

	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-t.closed[stock]:
				return

			case data := <-ch:
				tx, err := transaction.New(context.TODO(), nil)
				if err != nil {
					log.Fatal(err)
				}

				if err := t.broker.TransmitStockBatch(tx, stock, data); err != nil {
					log.Fatal(err)
				}
			}
		}
	}()

	return nil
}

func (t *Transmitter) UnsubscribeRelayer(stock string) error {
	delete(t.closed, stock)
	return nil
}
