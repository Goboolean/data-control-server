package adapter

import (
	"context"

	"github.com/Goboolean/stock-fetch-server/internal/infrastructure/postgresql"
	infra "github.com/Goboolean/stock-fetch-server/internal/infrastructure/transaction"
)

func newPsql(ctx context.Context) infra.Transactioner {
	instance := postgresql.NewInstance()
	tx, err := instance.Begin()

	if err != nil {
		panic(err)
	}

	return postgresql.NewTransaction(tx, ctx)
}
