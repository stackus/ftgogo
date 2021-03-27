package adapters

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func NewOrdersPlacedCounter() prometheus.Counter {
	return promauto.NewCounter(prometheus.CounterOpts{Name: "placed_orders"})
}
