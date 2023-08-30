package meta

import (
	"database/sql"

	"github.com/Goboolean/fetch-server.v1/internal/adapter/transaction"
	"github.com/Goboolean/fetch-server.v1/internal/domain/port"
	"github.com/Goboolean/fetch-server.v1/internal/domain/port/out"
	"github.com/Goboolean/fetch-server.v1/internal/domain/vo"
	"github.com/Goboolean/fetch-server.v1/internal/infrastructure/rdbms"
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

func (a *Adapter) GetStockMetadata(tx port.Transactioner, stockId string) (vo.StockAggsMeta, error) {
	q := a.rdbms.WithTx(tx.(*transaction.TxSession).P.Transaction().(*sql.Tx))

	dto, err := q.GetStockMeta(tx.Context(), stockId)

	return vo.StockAggsMeta{
		StockID: dto.ID,
	}, err
}

func (a *Adapter) GetAllStockMetadata(tx port.Transactioner) ([]vo.StockAggsMeta, error) {
	q := a.rdbms.WithTx(tx.(*transaction.TxSession).P.Transaction().(*sql.Tx))

	dtos, err := q.GetAllStockMetaList(tx.Context())
	if err != nil {
		return nil, err
	}

	metaList := make([]vo.StockAggsMeta, 0, len(dtos))

	for _, dto := range dtos {
		metaList = append(metaList, vo.StockAggsMeta{
			StockID: dto.ID,
		})
	}

	return metaList, nil
}
