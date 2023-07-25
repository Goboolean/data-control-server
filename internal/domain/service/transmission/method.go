package transmission

import (
	"context"
	"log"

	"github.com/Goboolean/fetch-server/internal/domain/entity"
)


func (t *Transmitter) SubscribeRelayer(ctx context.Context, stockId string) error {
	received := make([]*entity.StockAggregate, 0)

	if err := t.broker.CreateStockBroker(ctx, stockId); err != nil {
		return err
	}

	if err := t.s.StoreStock(stockId); err != nil {
		return err
	}

	ctx = t.s.Map[stockId].Context()

	ch, err := t.relayer.Subscribe(ctx, stockId)
	if err != nil {
		if err := t.s.UnstoreStock(stockId); err != nil {
			log.Println(err)
		}
		return err
	}


	go func(ctx context.Context) {
		for {
			select {
			case <- ctx.Done():
				return

			case data, ok := <-ch:

				if !ok {
					return
				}

				received = append(received, data)

				if len(received) % t.batchSize == 0 {
					ctx, cancel := context.WithCancel(ctx)
					defer cancel()

					if err := t.broker.TransmitStockBatch(ctx, stockId, received); err != nil {
						log.Println(err)

						continue
					}

					received = received[:0]
				}

			}
		}
	}(ctx)

	return nil
}


func (t *Transmitter) UnsubscribeRelayer(stockId string) error {
	return t.s.UnstoreStock(stockId)
}


func (t *Transmitter) IsStockTransmittable(stockId string) bool {
	return t.s.StockExists(stockId)
}
