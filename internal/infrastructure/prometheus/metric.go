package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)



func (s *Server) FetchCounter() func() prometheus.Counter {
	counter := promauto.NewCounter(prometheus.CounterOpts{
		Name: "StockCounter",
	})

	return func() prometheus.Counter {
		return counter
	}
}


func (s *Server) StoreCounter() func() prometheus.Counter {
	counter := promauto.NewCounter(prometheus.CounterOpts{
		Name: "StoreCounter",
	})

	return func() prometheus.Counter {
		return counter
	}
}


func (s *Server) BrokerCounter() func() prometheus.Counter {
	counter := promauto.NewCounter(prometheus.CounterOpts{
		Name: "BrokerCounter",
	})

	return func() prometheus.Counter {
		return counter
	}
}