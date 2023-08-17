package in

import "github.com/Goboolean/fetch-server/internal/domain/vo"

type RelayerPort interface {
	PlaceStockFormBatch([]*vo.StockAggregateForm)
}
