package stock

import (
	"database/sql"

	"github.com/Goboolean/stock-fetch-server/internal/adapter/transaction"
	"github.com/Goboolean/stock-fetch-server/internal/domain/port"
)

func (a *StockAdapter) CreateStoreLog(tx port.Transactioner, stock string) error {

	a.postgres.WithTx(tx.(*transaction.Transaction).P.Transaction().(*sql.Tx))

	if err := a.postgres.CreateAccessInfo(tx.Context()); err != nil {
		return err
	}

	return nil
}
