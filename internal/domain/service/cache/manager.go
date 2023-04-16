package cache

import (
	"context"

	adaptertx "github.com/Goboolean/data-control-server/internal/adapter/transaction"
)



func (m *StockCacheManager) SynchronizeDatabase(stock string ) {

	ctx := context.TODO()

	q := m.QMap[stock]

	q.Lock()
	defer q.Unlock()

	tx := adaptertx.NewTransaction(ctx)
	defer tx.Commit()

	if err := m.adapter.StoreStock(tx, stock, q.batch); err != nil {
		tx.Rollback()
	}

	if err := m.adapter.CreateStoreLog(tx, stock); err != nil {
		tx.Rollback()
	}
}