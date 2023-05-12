package receiver

import (
	"github.com/Goboolean/stock-fetch-server/internal/domain/value"
	"github.com/Goboolean/stock-fetch-server/internal/infrastructure/buycycle"
	"github.com/polygon-io/client-go/websocket/models"
)



func (s *StockReceiveAdapter) OnReceiveBuycycle(stock buycycle.StockAggregate) error {
	data := value.StockAggregateForm{}
	return s.port.PlaceStockFormBatch([]value.StockAggregateForm{data})
}



func (s *StockReceiveAdapter) OnReceivePolygonStockAggs(models.EquityAgg) error {

	data := value.StockAggregateForm{}
	return s.port.PlaceStockFormBatch([]value.StockAggregateForm{data})
}