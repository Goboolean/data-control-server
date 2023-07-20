package transmission

import (
	"context"
	"log"
)


func (t *Transmitter) SubscribeRelayer(ctx context.Context, stockId string) error {
	
	ch, err := t.relayer.Subscribe(stockId)
	if err != nil {
		return err
	}

	if err := t.s.StoreStock(stockId); err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <- t.s.Ctx.Done():
				return

			case data := <-ch:
				if err != nil {
					log.Fatal(err)
				}

				if err := t.broker.TransmitStockBatch(ctx, stockId, data); err != nil {
					log.Fatal(err)
				}
			}
		}
	}()

	return nil
}


func (t *Transmitter) UnsubscribeRelayer(ctx context.Context, stockId string) error {
	return t.s.UnstoreStock(stockId)
}


func (t *Transmitter) IsStockTransmittable(stockId string) bool {
	return t.s.StockExists(stockId)
}
