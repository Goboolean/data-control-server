package adaptertx

import (
	"context"

	"github.com/Goboolean/data-control-server/internal/infrastructure/postgresql"
	infratx "github.com/Goboolean/data-control-server/internal/infrastructure/transaction"
)



func NewPsqlTx(ctx context.Context) infratx.TransactionHandler {
	instance := postgresql.NewInstance()
	session, err := instance.Begin()

	if err != nil {
		panic(err)
	}

	return postgresql.NewTransaction(session, ctx);
}