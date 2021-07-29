package adapters

import (
	"github.com/stackus/ftgogo/order/internal/application/ports"
)

// TODO use in integration test of GRPC; Features?
var InmemCounters = map[string]*InmemCounter{}

type InmemCounter struct {
	count float64
}

var _ ports.Counter = (*InmemCounter)(nil)

func NewInmemCounter(name string) *InmemCounter {
	counter := &InmemCounter{count: 0}

	InmemCounters[name] = counter

	return counter
}

func (c *InmemCounter) Inc() {
	c.count += 1
}
