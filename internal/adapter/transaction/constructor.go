package transaction

import (
	"context"

	"github.com/Goboolean/shared/pkg/resolver"
)

type Transaction struct {
	M   resolver.Transactioner
	P   resolver.Transactioner
	R   resolver.Transactioner
	ctx context.Context
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

	return nil
}

func (t *Transaction) Context() context.Context {
	return t.ctx
}



type Option struct {
	Mongo    bool
	Postgres bool
	Redis    bool
}

func New(ctx context.Context, o *Option) (instance *Transaction, err error) {

	f := NewFactory()

	if o != nil {
		o = &Option{
			Mongo: true,
			Postgres: true,
			Redis: true,
		}	
	}

	if o.Mongo {
		instance.M, err = f.m.NewTx(ctx)
		if err != nil {
			return
		}
	}

	if o.Postgres {
		instance.P, err = f.p.NewTx(ctx)
		if err != nil {
			return
		}
	}

	if o.Redis {
		instance.R, err = f.r.NewTx(ctx)
		if err != nil {
			return
		}
	}

	return
}
