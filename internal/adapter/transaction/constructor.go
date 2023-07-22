package transaction

import (
	"context"
	"os"
	"sync"

	"github.com/Goboolean/fetch-server/internal/domain/port"
	"github.com/Goboolean/fetch-server/internal/infrastructure/redis"
	"github.com/Goboolean/shared/pkg/mongo"
	"github.com/Goboolean/shared/pkg/rdbms"
	"github.com/Goboolean/shared/pkg/resolver"
)




type Tx struct {
	m *mongo.DB
	p *rdbms.PSQL
	r *redis.Redis
}

var (
	once sync.Once
	instance *Tx
)


func New() port.TX {

	once.Do(func() {
		instance = &Tx{
			r: redis.NewInstance(&resolver.ConfigMap{
				"HOST":     os.Getenv("REDIS_HOST"),
				"PORT":     os.Getenv("REDIS_PORT"),
				"USER":     os.Getenv("REDIS_USER"),
				"PASSWORD": os.Getenv("REDIS_PASS"),
			}),
	
			m: mongo.NewDB(&resolver.ConfigMap{
				"HOST":     os.Getenv("MONGO_HOST"),
				"PORT":     os.Getenv("MONGO_PORT"),
				"USER":     os.Getenv("MONGO_USER"),
				"PASSWORD": os.Getenv("MONGO_PASS"),
			}),
	
			p: rdbms.NewDB(&resolver.ConfigMap{
				"HOST":     os.Getenv("PSQL_HOST"),
				"PORT":     os.Getenv("PSQL_PORT"),
				"USER":     os.Getenv("PSQL_USER"),
				"PASSWORD": os.Getenv("PSQL_PASS"),
			}),
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

	r, err := t.r.NewTx(ctx)
	if err != nil {
		return nil, err
	}

	var a *TxSession = &TxSession{M: m, P: p, R: r}

	return a, nil
}
