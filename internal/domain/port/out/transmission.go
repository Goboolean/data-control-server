package out

import (
	"github.com/Goboolean/fetch-server/internal/domain/entity"
)

type TransmissionPort interface {
	TransmitStockBatch(string, []*entity.StockAggregate) error
}
