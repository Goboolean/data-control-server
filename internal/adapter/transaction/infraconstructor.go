package transaction

import (
	"context"

	"github.com/Goboolean/shared-packages/pkg/mongo"
	"github.com/Goboolean/shared-packages/pkg/rdbms"
	"github.com/Goboolean/shared-packages/pkg/resolver"
	"github.com/Goboolean/stock-fetch-server/internal/infrastructure/rediscache"
)



func newM(ctx context.Context, f *Factory) (resolver.Transactioner, error) {
	session, err := f.m.StartSession()
	if err != nil {
		return nil, err
	}

	return mongo.NewTransaction(session, ctx), nil
}



func newP(ctx context.Context, f *Factory) (resolver.Transactioner, error) {
	tx, err := f.p.Begin()
	if err != nil {
		return nil, err
	}

	return rdbms.NewTransaction(tx, ctx), nil
}



func newR(ctx context.Context, f *Factory) (resolver.Transactioner, error) {
	var pipe = f.r.TxPipeline()
	return rediscache.NewTransaction(pipe, ctx), nil
}



func newK(ctx context.Context, f *Factory) (resolver.Transactioner, error) {
	return nil, nil
}