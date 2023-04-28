package persistence

import (
	"context"

	"github.com/Goboolean/stock-fetch-server/internal/adapter/transaction"
)

func (m *PersistenceManager) SynchronizeDatabase(ctx context.Context, stock string) error {

	tx := adapter.NewMongoRdbmsRedis(ctx)

	batch, err := m.db.EmptyCache(tx, stock)

	if err != nil {
		tx.Rollback()
		return err
	}

	if err := m.db.StoreStock(tx, stock, batch); err != nil {
		tx.Rollback()
		return err
	}

	if err := m.db.CreateStoreLog(tx, stock); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
