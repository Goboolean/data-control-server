package meta

import (
	"database/sql"

	"github.com/Goboolean/fetch-server/internal/adapter/transaction"
	"github.com/Goboolean/fetch-server/internal/domain/entity"
	"github.com/Goboolean/fetch-server/internal/domain/port"
	"github.com/Goboolean/shared/pkg/rdbms"
)



type Adapter struct {
	q *rdbms.Queries
}

func NewAdapter(q *rdbms.Queries) *Adapter {
	return &Adapter{q: q}
}


func (a *Adapter) CheckStockExists(tx port.Transactioner, stockId string) (bool, error) {
	q := a.q.WithTx(tx.(*transaction.TxSession).R.Transaction().(*sql.Tx))
	return q.CheckStockExist(tx.Context(), stockId)
}


func (a *Adapter) GetStockMetadata(tx port.Transactioner, stockId string) (entity.StockAggsMeta, error) {
	q := a.q.WithTx(tx.(*transaction.TxSession).R.Transaction().(*sql.Tx))

	dto, err := q.GetStockMeta(tx.Context(), stockId)

	return entity.StockAggsMeta{
		StockID:   dto.ID,
	}, err
}


func (a *Adapter) GetAllStockMetadata(tx port.Transactioner) ([]entity.StockAggsMeta, error) {
	q := a.q.WithTx(tx.(*transaction.TxSession).R.Transaction().(*sql.Tx))

	dtos, err := q.GetAllStockMetaList(tx.Context())
	if err != nil {
		return nil, err
	}

	metaList := make([]entity.StockAggsMeta, 0, len(dtos))

	for _, dto := range dtos {
		metaList = append(metaList, entity.StockAggsMeta{
			StockID:   dto.ID,
		})
	}

	return metaList, nil
}
