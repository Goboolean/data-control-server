package persistence

import (
	"context"
	"log"

	adapter "github.com/Goboolean/stock-fetch-server/internal/adapter/transaction"
)

func (m *PersistenceManager) SubscribeRelayer(stock string) error {
	ch, err := m.relayer.Subscribe(stock)

	if err != nil {
		return err
	}

	go func() {

		for {
			select {
			case <-m.running[stock]:
				return

			case data := <-ch:
				tx := adapter.NewRedis(context.TODO())
				if err := m.db.InsertOnCache(tx, stock, data); err != nil {
					log.Fatal(err)
				}
			}
		}
	}()

	return nil
}

func (m *PersistenceManager) UnsubscribeRelayer(stock string) error {
	delete(m.running, stock)
	return nil
}
