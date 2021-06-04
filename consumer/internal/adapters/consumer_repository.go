package adapters

import (
	"github.com/stackus/edat/es"
)

type ConsumerRepository interface {
	es.AggregateRepository
}
