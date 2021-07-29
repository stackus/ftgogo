package adapters

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func NewPrometheusCounter(name string) prometheus.Counter {
	return promauto.NewCounter(prometheus.CounterOpts{Name: name})
}
