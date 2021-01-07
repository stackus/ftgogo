package adapters

import (
	"github.com/stackus/edat/es"
	"github.com/stackus/ftgogo/kitchen/internal/domain"
)

func NewTicketRepository(store es.AggregateRootStore) domain.TicketRepository {
	return es.NewAggregateRootRepository(domain.NewTicket, store)
}
