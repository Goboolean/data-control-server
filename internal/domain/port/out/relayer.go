package out

import "github.com/Goboolean/fetch-server/internal/domain/entity"



type RelayerPort interface {
	FetchStock(stockId string, stockMeta entity.StockAggsMeta) error
	StopFetchingStock(stockId string) error
}