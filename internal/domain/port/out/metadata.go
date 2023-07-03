package out

import (
	"github.com/Goboolean/fetch-server/internal/domain/port"
	"github.com/Goboolean/fetch-server/internal/domain/value"
)

type StockMetadataPort interface {
	GetStockType(port.Transactioner, string) (value.StockType, error)
	StockExists(port.Transactioner, string) (bool, error)
}
