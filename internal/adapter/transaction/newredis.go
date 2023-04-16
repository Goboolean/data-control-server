package adaptertx

import (
	"context"

	"github.com/Goboolean/data-control-server/internal/infrastructure/redis"
	infratx "github.com/Goboolean/data-control-server/internal/infrastructure/transaction"
)


func NewRedisTx(ctx context.Context) infratx.TransactionHandler {
	instance := redis.New()
	pipe := instance.TxPipeline()

	return redis.NewTransaction(pipe, ctx)
}