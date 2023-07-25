package transmission

import (
	"context"
	"errors"
	"log"

	"github.com/Goboolean/fetch-server/internal/domain/entity"
)


func (t *Transmitter) SubscribeRelayer(ctx context.Context, stockId string) error {
	received := make([]*entity.StockAggregate, 0)

	if exists := t.s.StockExists(stockId); exists {
		return errors.New("stock already exists")
	}

	ch, err := t.relayer.Subscribe(ctx, stockId)
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

			case data, ok := <-ch:

				if !ok {
					return
				}

				received = append(received, data)

				if len(received) % t.batchSize == 0 {
					ctx, cancel := context.WithCancel(t.s.Ctx.Context())
					defer cancel()

					if err := t.broker.TransmitStockBatch(ctx, stockId, received); err != nil {
						log.Println(err)

						continue
					}

					received = received[:0]
				}

			}
		}
	}()

	return nil
}


func (t *Transmitter) UnsubscribeRelayer(stockId string) error {
	return t.s.UnstoreStock(stockId)
}


func (t *Transmitter) IsStockTransmittable(stockId string) bool {
	return t.s.StockExists(stockId)
}
