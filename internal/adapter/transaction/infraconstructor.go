package transaction

import (
	"context"

	"github.com/Goboolean/shared-packages/pkg/mongo"
	"github.com/Goboolean/shared-packages/pkg/rdbms"
	"github.com/Goboolean/shared-packages/pkg/resolver"
	"github.com/Goboolean/stock-fetch-server/internal/infrastructure/rediscache"
)



func newMongo(ctx context.Context) (resolver.Transactioner, error) {
	instance := mongo.New()
	session, err := instance.StartSession()

	if err != nil {
		return err, nil
	}

	return mongo.NewTransaction(session, ctx), nil
}



func newPSQL(ctx context.Context) (resolver.Transactioner, error) {
	instance := rdbms.New()
	tx, err := instance.Begin()

	if err != nil {
		return err, nil
	}

	return rdbms.NewTransaction(tx, ctx), nil
}



func newRedis(ctx context.Context) (resolver.Transactioner, error) {
	instance := rediscache.NewInstance()
	var pipe = instance.TxPipeline()

	return rediscache.NewTransaction(pipe, ctx), nil
}



func newKafka(ctx context.Context) (resolver.Transactioner, error) {
	return nil, nil
}