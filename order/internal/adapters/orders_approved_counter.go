package adapters

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func NewOrdersApprovedCounter() prometheus.Counter {
	return promauto.NewCounter(prometheus.CounterOpts{Name: "approved_orders"})
}
