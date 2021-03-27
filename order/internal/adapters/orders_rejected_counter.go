package adapters

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func NewOrdersRejectedCounter() prometheus.Counter {
	return promauto.NewCounter(prometheus.CounterOpts{Name: "rejected_orders"})
}
