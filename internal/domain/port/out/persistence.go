package out

import (
	"github.com/Goboolean/fetch-server/internal/domain/port"
	"github.com/Goboolean/fetch-server/internal/domain/entity"
)

type StockPersistencePort interface {
	EmptyCache(port.Transactioner, string) ([]*entity.StockAggregate, error)
	StoreStock(port.Transactioner, string, []*entity.StockAggregate) error
	CreateStoreLog(port.Transactioner, string) error
	InsertOnCache(port.Transactioner, string, []*entity.StockAggregate) error
}
