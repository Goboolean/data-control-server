package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)



var (
	StockCounter prometheus.Counter
	StoreCounter prometheus.Counter
	MQCounter prometheus.Counter
	RequestCounter prometheus.Counter
)



func init() {

	StockCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "StockCounter",
	})

	StoreCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "StoreCounter",
	})

	MQCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "MQCounter",
	})

	RequestCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "RequestCounter",
	})	
}