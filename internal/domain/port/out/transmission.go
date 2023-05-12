package out

import (
	"github.com/Goboolean/stock-fetch-server/internal/domain/port"
	"github.com/Goboolean/stock-fetch-server/internal/domain/value"
)



type TransmissionPort interface {
	TransmitStockBatch(port.Transactioner, string, []value.StockAggregate) error
}