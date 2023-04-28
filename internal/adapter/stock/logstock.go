package stock

import (
	"github.com/Goboolean/stock-fetch-server/internal/adapter/transaction"
	"github.com/Goboolean/stock-fetch-server/internal/domain/port"
	"github.com/Goboolean/stock-fetch-server/internal/infrastructure/postgresql"
)

func (a *StockAdapter) CreateStoreLog(tx port.Transactioner, stock string) error {

	q := postgresql.SetTx(tx.(*adapter.Transaction).Psql.Transaction().(*postgresql.Transaction))

	if err := q.CreateAccessInfo(tx.Context()); err != nil {
		return err
	}

	return nil
}
