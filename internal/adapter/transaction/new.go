package adapter

import (
	"context"

	"github.com/Goboolean/stock-fetch-server/internal/domain/port"
)

func NewMongoRdbmsRedis(ctx context.Context) port.Transactioner {

	return &Transaction{
		ctx: ctx,

		Mongo: newMongo(ctx),
		Psql:  newPsql(ctx),
	}
}

func NewRedis(ctx context.Context) port.Transactioner {

	return &Transaction{
		ctx:   ctx,
		Redis: newRedis(ctx),
	}
}
