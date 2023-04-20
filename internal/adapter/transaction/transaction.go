package adaptertx

import (
	"context"

	"github.com/Goboolean/data-control-server/internal/domain/port"
	"github.com/Goboolean/data-control-server/internal/infrastructure/transaction"
)



type Transaction struct {
	MongoTx infratx.TransactionHandler
	PsqlTx  infratx.TransactionHandler
	RedisTx infratx.TransactionHandler
	ctx context.Context
}



func (t *Transaction) Commit() error {

	if err := t.MongoTx.Commit(); err != nil {
		return err
	}

	if err := t.PsqlTx.Commit(); err != nil {
		return err
	}

	if err := t.RedisTx.Commit(); err != nil {
		return err
	}

	return nil
}



func (t *Transaction) Rollback() error {

	if err := t.MongoTx.Rollback(); err != nil {
		return err
	}

	if err := t.PsqlTx.Rollback(); err != nil {
		return err
	}

	if err := t.RedisTx.Rollback(); err != nil {
		return err
	}

	return nil
}

func (t *Transaction) Context() context.Context {
	return t.ctx
}

func NewTransaction(ctx context.Context) port.Transactioner {

	return &Transaction{
		ctx: ctx,

		MongoTx: NewMongoTx(ctx),
		PsqlTx: NewPsqlTx(ctx),
		RedisTx: NewRedisTx(ctx),
	}
}

