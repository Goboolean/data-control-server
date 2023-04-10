package adaptertx

import (
	"context"

	"github.com/Goboolean/data-control-server/internal/infrastructure/mongodb"
	infratx "github.com/Goboolean/data-control-server/internal/infrastructure/transaction"
)




func NewMongoTx(ctx context.Context) infratx.TransactionHandler {
	instance := mongodb.NewInstance()
	session, err := instance.StartSession()

	if err != nil {
		panic(err)
	}

	return mongodb.NewTransaction(&session, ctx);
}