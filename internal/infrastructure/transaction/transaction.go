package infratx

import "context"

type TransactionHandler interface {
	Commit() error
	Rollback() error
	Context() context.Context
	Transaction() interface{}
}
