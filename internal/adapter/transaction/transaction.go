package adapter

import (
	"context"

	infra "github.com/Goboolean/data-control-server/internal/infrastructure/transaction"
)

type Transaction struct {
	Mongo infra.Transactioner
	Psql  infra.Transactioner
	Redis infra.Transactioner
	ctx   context.Context
}

func (t *Transaction) Commit() error {

	if err := t.Mongo.Commit(); err != nil {
		return err
	}

	if err := t.Psql.Commit(); err != nil {
		return err
	}

	if err := t.Redis.Commit(); err != nil {
		return err
	}

	return nil
}

func (t *Transaction) Rollback() error {

	if err := t.Mongo.Rollback(); err != nil {
		return err
	}

	if err := t.Psql.Rollback(); err != nil {
		return err
	}

	if err := t.Redis.Rollback(); err != nil {
		return err
	}

	return nil
}

func (t *Transaction) Context() context.Context {
	return t.ctx
}
