package transaction

import (
	"context"
	"sync"

	"github.com/Goboolean/fetch-server/internal/domain/port"
	"github.com/Goboolean/shared/pkg/mongo"
	"github.com/Goboolean/shared/pkg/rdbms"
)




type Tx struct {
	m *mongo.DB
	p *rdbms.PSQL
}

var (
	once sync.Once
	instance *Tx
)


func New(m *mongo.DB, p *rdbms.PSQL) port.TX {

	once.Do(func() {
		instance = &Tx{
			m: m,
			p: p,
		}
	})

	return instance
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

	var a *TxSession = &TxSession{M: m, P: p}

	return a, nil
}
