package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)




var (
	FetchCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "FetchCounter",
	})
	StoreCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "StoreCounter",
	})
	BrokerCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "BrokerCounter",
	})
)
