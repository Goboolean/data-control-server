package adaptertx

import (
	"context"
	"github.com/Goboolean/data-control-server/internal/infrastructure/transaction"
)



type Transaction struct {
	MongoTx infratx.TransactionHandler
	PsqlTx  infratx.TransactionHandler
	RedisTx infratx.TransactionHandler
	ctx context.Context
}

func (t *Transaction) Commit() {
	if err := t.MongoTx.Commit(); err != nil {
		panic(err)
	}
	if err := t.PsqlTx.Commit(); err != nil {
		panic(err)
	}
	if err := t.RedisTx.Commit(); err != nil {
		panic(err)
	}
}

func (t *Transaction) Rollback() {
	if err := t.MongoTx.Rollback(); err != nil {
		panic(err)
	}
	if err := t.PsqlTx.Rollback(); err != nil {
		panic(err)
	}
	if err := t.RedisTx.Rollback(); err != nil {
		panic(err)
	}
}

func (t *Transaction) Context() context.Context {
	return t.ctx
}

func NewTransaction() *Transaction {
	ctx := context.Background()

	return &Transaction{
		ctx: ctx,

		MongoTx: NewMongoTx(ctx),
		PsqlTx: NewPsqlTx(ctx),
		RedisTx: NewRedisTx(ctx),
	}
}

