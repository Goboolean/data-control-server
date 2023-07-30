package port

import "context"


type Transactioner interface {
	Commit() error
	Rollback() error
	Context() context.Context
}

type TX interface {
	Transaction(context.Context) (Transactioner, error)
}