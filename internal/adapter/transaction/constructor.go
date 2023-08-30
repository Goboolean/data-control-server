package transaction

import (
	"context"

	"github.com/Goboolean/fetch-server.v1/internal/domain/port"
	"github.com/Goboolean/fetch-server.v1/internal/infrastructure/mongo"
	"github.com/Goboolean/fetch-server.v1/internal/infrastructure/rdbms"
)

type Tx struct {
	m *mongo.DB
	p *rdbms.PSQL
}

func New(m *mongo.DB, p *rdbms.PSQL) port.TX {

	return &Tx{
		m: m,
		p: p,
	}
}

func (t *Tx) Transaction(ctx context.Context) (port.Transactioner, error) {

	m, err := t.m.NewTx(ctx)
	if err != nil {
		return nil, err
	}

	p, err := t.p.NewTx(ctx)
	if err != nil {
		return nil, err
	}

	return &TxSession{M: m, P: p, ctx: ctx}, nil
}
