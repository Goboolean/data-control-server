package out

import "github.com/Goboolean/stock-fetch-server/internal/domain/value"


type StockMetadataPort interface {
	GetStockType(string) (value.StockType, error)
	StockExists(string) (bool, error)
}