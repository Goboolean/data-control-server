package mongo

import (
	"context"

	"github.com/Goboolean/shared/pkg/resolver"
	"go.mongodb.org/mongo-driver/mongo"
)

type Transaction struct {
	session mongo.Session
	ctx     context.Context
}

func (d *Transaction) Commit() error {
	return d.session.CommitTransaction(d.ctx)
}

func (d *Transaction) Rollback() error {
	defer d.session.EndSession(d.ctx)
	return d.session.AbortTransaction(d.ctx);
}

func (d *Transaction) Context() context.Context {
	return d.ctx
}

func (d *Transaction) Transaction() interface{} {
	return d.session
}

func NewTransaction(session mongo.Session, ctx context.Context) resolver.Transactioner {
	return &Transaction{session: session, ctx: ctx}
}
