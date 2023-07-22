package out

import (
	"github.com/Goboolean/fetch-server/internal/domain/port"
	"github.com/Goboolean/fetch-server/internal/domain/entity"
)

type StockMetadataPort interface {
	CheckStockExists(port.Transactioner, string) (bool, error)
	GetStockMetadata(port.Transactioner, string) (entity.StockAggsMeta, error)
	GetAllStockMetadata(port.Transactioner) ([]entity.StockAggsMeta, error)
}
