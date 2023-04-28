package adapter

import (
	"context"

	"github.com/Goboolean/stock-fetch-server/internal/infrastructure/mongodb"
	infra "github.com/Goboolean/stock-fetch-server/internal/infrastructure/transaction"
)

func newMongo(ctx context.Context) infra.Transactioner {
	instance := mongodb.NewInstance()
	session, err := instance.StartSession()

	if err != nil {
		panic(err)
	}

	return mongodb.NewTransaction(session, ctx)
}
