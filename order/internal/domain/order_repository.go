package domain

import (
	"github.com/stackus/edat/es"
)

type OrderRepository interface {
	es.AggregateRepository
}
