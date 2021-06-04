package adapters

import (
	"github.com/stackus/edat/es"

	"github.com/stackus/ftgogo/kitchen/internal/domain"
)

func NewTicketAggregateRootRepository(store es.AggregateRootStore) TicketRepository {
	return es.NewAggregateRootRepository(domain.NewTicket, store)
}
