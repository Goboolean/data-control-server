package infra

import "context"

type Transactioner interface {
	Commit() error
	Rollback() error
	Context() context.Context
	Transaction() interface{}
}
