package transaction

import (
	"context"

	"github.com/Goboolean/shared-packages/pkg/mongo"
	"github.com/Goboolean/shared-packages/pkg/rdbms"
	"github.com/Goboolean/shared-packages/pkg/resolver"
	"github.com/Goboolean/stock-fetch-server/internal/infrastructure/rediscache"
)



func newMongo(ctx context.Context) resolver.Transactioner {
	instance := mongo.New()
	session, err := instance.StartSession()

	if err != nil {
		panic(err)
	}

	return mongo.NewTransaction(session, ctx)
}



func newPSQL(ctx context.Context) resolver.Transactioner {
	instance := rdbms.New()
	tx, err := instance.Begin()

	if err != nil {
		panic(err)
	}

	return rdbms.NewTransaction(tx, ctx)
}



func newRedis(ctx context.Context) resolver.Transactioner {
	instance := rediscache.NewInstance()
	var pipe = instance.TxPipeline()

	return rediscache.NewTransaction(pipe, ctx)
}



func newKafka(ctx context.Context) resolver.Transactioner {
	return nil
}