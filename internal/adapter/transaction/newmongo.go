package adapter

import (
	"context"

	"github.com/Goboolean/data-control-server/internal/infrastructure/mongodb"
	infra "github.com/Goboolean/data-control-server/internal/infrastructure/transaction"
)

func newMongo(ctx context.Context) infra.Transactioner {
	instance := mongodb.NewInstance()
	session, err := instance.StartSession()

	if err != nil {
		panic(err)
	}

	return mongodb.NewTransaction(session, ctx)
}
