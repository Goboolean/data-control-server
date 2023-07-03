package persistence

import (
	"context"
	"log"

	"github.com/Goboolean/fetch-server/internal/adapter/transaction"
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
				tx, err := transaction.New(context.TODO(), &transaction.Option{
					Redis: true,
				})
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
