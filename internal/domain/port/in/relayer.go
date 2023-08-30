package in

import "github.com/Goboolean/fetch-server.v1/internal/domain/vo"

type RelayerPort interface {
	PlaceStockFormBatch([]*vo.StockAggregateForm)
}
