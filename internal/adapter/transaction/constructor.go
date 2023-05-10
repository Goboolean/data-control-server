package transaction

import (
	"context"

	"github.com/Goboolean/shared-packages/pkg/resolver"
)

type Transaction struct {
	M resolver.Transactioner
	P resolver.Transactioner
	R resolver.Transactioner
	K resolver.Transactioner
	ctx   context.Context
}

func (t *Transaction) Commit() error {

	if t.M != nil {
		if err := t.M.Commit(); err != nil {
			return err
		}
	}

	if t.P != nil {
		if err := t.P.Commit(); err != nil {
			return err
		}
	}

	if t.R != nil {
		if err := t.R.Commit(); err != nil {
			return err
		}
	}

	if t.R != nil {
		if err := t.K.Commit(); err != nil {
			return err
		}
	}

	return nil
}



func (t *Transaction) Rollback() error {

	if t.M != nil {
		if err := t.M.Rollback(); err != nil {
			return err
		}
	}

	if t.P != nil {
		if err := t.P.Rollback(); err != nil {
			return err
		}
	}

	if t.R != nil {
		if err := t.R.Rollback(); err != nil {
			return err
		}
	}

	if t.R != nil {
		if err := t.K.Rollback(); err != nil {
			return err
		}
	}

	return nil
}



func (t *Transaction) Context() context.Context {
	return t.ctx
}



type Option struct {
	Mongo    bool
	Postgres bool
	Redis    bool
	Kafka    bool
}

func New(ctx context.Context, o *Option) *Transaction {
	
	instance := &Transaction{}
	
	if o.Mongo {
		instance.M = newMongo(ctx)
	}

	if o.Postgres {
		instance.P = newPSQL(ctx)
	}

	if o.Redis {
		instance.R = newRedis(ctx)
	}

	if o.Kafka {
		instance.K = newKafka(ctx)
	}

	return instance
}