package meta

import (
	"database/sql"

	"github.com/Goboolean/fetch-server/internal/adapter/transaction"
	"github.com/Goboolean/fetch-server/internal/domain/entity"
	"github.com/Goboolean/fetch-server/internal/domain/port"
	"github.com/Goboolean/fetch-server/internal/domain/port/out"
	"github.com/Goboolean/shared/pkg/rdbms"
)



type Adapter struct {
	rdbms *rdbms.Queries
}

func NewAdapter(rdbms *rdbms.Queries) out.StockMetadataPort {
	return &Adapter{rdbms: rdbms}
}


func (a *Adapter) CheckStockExists(tx port.Transactioner, stockId string) (bool, error) {
	q := a.rdbms.WithTx(tx.(*transaction.TxSession).P.Transaction().(*sql.Tx))
	return q.CheckStockExist(tx.Context(), stockId)
}


func (a *Adapter) GetStockMetadata(tx port.Transactioner, stockId string) (entity.StockAggsMeta, error) {
	q := a.rdbms.WithTx(tx.(*transaction.TxSession).P.Transaction().(*sql.Tx))

	dto, err := q.GetStockMeta(tx.Context(), stockId)

	return entity.StockAggsMeta{
		StockID:   dto.ID,
	}, err
}


func (a *Adapter) GetAllStockMetadata(tx port.Transactioner) ([]entity.StockAggsMeta, error) {
	q := a.rdbms.WithTx(tx.(*transaction.TxSession).P.Transaction().(*sql.Tx))

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
