package receiver

import (
	"github.com/Goboolean/fetch-server/internal/domain/value"
	"github.com/Goboolean/fetch-server/internal/infrastructure/buycycle"
	"github.com/Goboolean/fetch-server/internal/infrastructure/prometheus"
	"github.com/polygon-io/client-go/websocket/models"
)

func (s *StockReceiveAdapter) OnReceiveBuycycle(batch []buycycle.StockAggregate) error {

	prometheus.DomesticStockCounter.Add(float64(len(batch)))

	data := value.StockAggregateForm{}
	return s.port.PlaceStockFormBatch([]value.StockAggregateForm{data})
}

func (s *StockReceiveAdapter) OnReceivePolygonStockAggs(batch []models.EquityAgg) error {

	prometheus.ForeignStockCounter.Add(float64(len(batch)))

	data := value.StockAggregateForm{}
	return s.port.PlaceStockFormBatch([]value.StockAggregateForm{data})
}
