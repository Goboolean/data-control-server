package in

import "github.com/Goboolean/fetch-server/internal/domain/value"

type RelayerPort interface {
	PlaceStockFormBatch([]value.StockAggregateForm) error
}
