package adapter

import (
	"context"

	"github.com/Goboolean/stock-fetch-server/internal/infrastructure/rediscache"
	infra "github.com/Goboolean/stock-fetch-server/internal/infrastructure/transaction"
)

func newRedis(ctx context.Context) infra.Transactioner {
	instance := rediscache.NewInstance()
	var pipe = instance.TxPipeline()

	return rediscache.NewTransaction(pipe, ctx)
}
