package outport

import (
	"github.com/Goboolean/stock-fetch-server/internal/domain/port"
	"github.com/Goboolean/stock-fetch-server/internal/domain/value"
)

type StockPersistencePort interface {
	EmptyCache(port.Transactioner, string) ([]value.StockAggregate, error)
	StoreStock(port.Transactioner, string, []value.StockAggregate) error
	CreateStoreLog(port.Transactioner, string) error
	InsertOnCache(port.Transactioner, string, []value.StockAggregate) error
}
