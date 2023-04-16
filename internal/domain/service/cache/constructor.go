package cache

import (
	"github.com/Goboolean/data-control-server/internal/adapter/stock"
	outport "github.com/Goboolean/data-control-server/internal/domain/port/out"
)


type StockCacheManager struct {
	QMap map[string]*stockQueue

	adapter outport.StockPort
}

var instance *StockCacheManager



func NewStockCacheManager() *StockCacheManager {

	if instance == nil {
		instance = &StockCacheManager{
			QMap: make(map[string]*stockQueue),
			adapter: stock.NewStockAdapter(),
		}
	}

	return instance
}