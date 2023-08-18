package transaction

import (
	"context"

	"github.com/Goboolean/shared/pkg/resolver"
)

type TxSession struct {
	M   resolver.Transactioner
	P   resolver.Transactioner
	ctx context.Context
}


func (t *TxSession) Commit() error {

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

	return nil
}


func (t *TxSession) Rollback() error {

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

	return nil
}

func (t *TxSession) Context() context.Context {
	return t.ctx
}
