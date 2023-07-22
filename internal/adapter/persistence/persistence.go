package persistence

import (
	"github.com/Goboolean/fetch-server/internal/domain/entity"
	"github.com/Goboolean/fetch-server/internal/domain/port"
	"github.com/Goboolean/fetch-server/internal/domain/port/out"
	"github.com/Goboolean/shared/pkg/mongo"
	"github.com/Goboolean/shared/pkg/rdbms"
)




type Adapter struct {
	rdbms *rdbms.Queries
	mongo *mongo.Queries
}

func NewAdapter(rdbms *rdbms.Queries, mongo *mongo.Queries) out.StockPersistencePort {
	return &Adapter{
		rdbms: rdbms,
		mongo: mongo,
	}
}


func (a *Adapter) EmptyCache(port.Transactioner, string) ([]*entity.StockAggregate, error)

func (a *Adapter) StoreStock(port.Transactioner, string, []*entity.StockAggregate) error


func (a *Adapter) CreateStoreLog(port.Transactioner, string) error
func (a *Adapter) InsertOnCache(port.Transactioner, string, []*entity.StockAggregate) error

