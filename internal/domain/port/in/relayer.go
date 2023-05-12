package in

import "github.com/Goboolean/stock-fetch-server/internal/domain/value"

type RelayerPort interface {
	PlaceStockFormBatch([]value.StockAggregateForm) error
}
