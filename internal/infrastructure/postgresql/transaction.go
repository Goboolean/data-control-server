package postgresql

import (
	"context"
	"github.com/Goboolean/data-control-server/internal/infrastructure/transaction"
	"database/sql"
)

type Transaction struct {
	tx  *sql.Tx
	ctx context.Context
}

func (d *Transaction) Commit() error {
	return d.tx.Commit()
}

func (d *Transaction) Rollback() error {
	return d.tx.Rollback()
}

func (d *Transaction) Context() context.Context {
	return d.ctx
}

func (d *Transaction) Transaction() interface{} {
	return d.tx
}

func NewTransaction(tx *sql.Tx, ctx context.Context) infra.Transactioner {
	return &Transaction{tx: tx, ctx: ctx}
}
