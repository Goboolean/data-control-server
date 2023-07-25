package transaction

import (
	"context"

	"github.com/Goboolean/fetch-server/internal/domain/port"
)



type MockTx struct {}

func NewMock() port.TX {
	return &MockTx{}
}

func (t *MockTx) Transaction(ctx context.Context) (port.Transactioner, error) {
	return &MockTxSession{ctx: ctx}, nil
}

type MockTxSession struct {
	ctx context.Context
}

func (t *MockTxSession) Commit() error {
	return nil
}

func (t *MockTxSession) Rollback() error {
	return nil
}

func (t *MockTxSession) Context() context.Context {
	return t.ctx
}