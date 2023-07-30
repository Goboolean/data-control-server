package in

import "github.com/Goboolean/fetch-server/internal/domain/entity"

type RelayerPort interface {
	PlaceStockFormBatch([]*entity.StockAggregateForm)
}
