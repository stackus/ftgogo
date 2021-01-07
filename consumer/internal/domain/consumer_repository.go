package domain

import (
	"github.com/stackus/edat/es"
)

type ConsumerRepository interface {
	es.AggregateRepository
}
