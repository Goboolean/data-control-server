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
