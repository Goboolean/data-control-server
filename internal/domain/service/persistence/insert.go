package persistence

import (
	"context"
	"log"
)

func (m *PersistenceManager) SubscribeRelayer(stock string) error {

	ch, err := m.relayer.Subscribe(stock)

	if err != nil {
		return err
	}

	go func() {

		for {
			select {
			case <-m.closed[stock]:
				return

			case data := <-ch:
				ctx := context.Background()

				tx, err := m.tx.Transaction(ctx)
				if err != nil {
					log.Fatal(err)
				}

				if err := m.db.InsertOnCache(tx, stock, data); err != nil {
					log.Fatal(err)
				}
			}
		}
	}()

	return nil
}

func (m *PersistenceManager) UnsubscribeRelayer(stock string) error {
	delete(m.closed, stock)
	return nil
}
