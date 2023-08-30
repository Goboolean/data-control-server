package out

import (
	"github.com/Goboolean/fetch-server.v1/internal/domain/port"
	"github.com/Goboolean/fetch-server.v1/internal/domain/vo"
)

type StockMetadataPort interface {
	CheckStockExists(port.Transactioner, string) (bool, error)
	GetStockMetadata(port.Transactioner, string) (vo.StockAggsMeta, error)
	GetAllStockMetadata(port.Transactioner) ([]vo.StockAggsMeta, error)
}
