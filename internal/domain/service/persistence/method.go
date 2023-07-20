package persistence

import (
	"context"
	"log"
)

func (m *PersistenceManager) SubscribeRelayer(ctx context.Context, stockId string) error {

	ch, err := m.relayer.Subscribe(stockId)
	if err != nil {
		return err
	}

	if err := m.s.StoreStock(stockId); err != nil {
		return err
	}

	go func() {

		for {
			select {
			case <- m.s.Map[stockId].Done():
				return

			case data := <-ch:

				tx, err := m.tx.Transaction(ctx)
				if err != nil {
					log.Fatal(err)
				}

				if err := m.db.InsertOnCache(tx, stockId, data); err != nil {
					log.Fatal(err)
				}
			}
		}
	}()

	return nil
}


func (m *PersistenceManager) UnsubscribeRelayer(ctx context.Context, stockId string) error {
	return m.s.UnstoreStock(stockId)
}


func (m *PersistenceManager) IsStockStoreable(stockId string) bool {
	return m.s.StockExists(stockId)
}


func (m *PersistenceManager) SynchronizeDatabase(ctx context.Context, stock string) error {

	tx, err := m.tx.Transaction(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	batch, err := m.db.EmptyCache(tx, stock)
	if err != nil {
		return err
	}

	if err := m.db.StoreStock(tx, stock, batch); err != nil {
		return err
	}

	if err := m.db.CreateStoreLog(tx, stock); err != nil {
		return err
	}

	return tx.Commit()
}
