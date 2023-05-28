package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)



var (
	DomesticStockCounter prometheus.Counter
	ForeignStockCounter prometheus.Counter
	DBCounter prometheus.Counter
	MQCounter prometheus.Counter	
)



func init() {

	DomesticStockCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "counts domestic stock received",
	})

	ForeignStockCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "counts foreign stock received",
	})

	DBCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "counts stock stored on db",
	})

	MQCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "counts stock sended to kafka",
	})
}