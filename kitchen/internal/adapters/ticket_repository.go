package adapters

import (
	"github.com/stackus/edat/es"
)

type TicketRepository interface {
	es.AggregateRepository
}
