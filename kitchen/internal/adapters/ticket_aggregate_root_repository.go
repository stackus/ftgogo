package adapters

import (
	"context"

	"github.com/stackus/edat/es"
	"github.com/stackus/errors"

	"github.com/stackus/ftgogo/kitchen/internal/domain"
)

type TicketAggregateRootRepository struct {
	es.AggregateRepository
}

func NewTicketAggregateRootRepository(store es.AggregateRootStore) *TicketAggregateRootRepository {
	return &TicketAggregateRootRepository{es.NewAggregateRootRepository(domain.NewTicket, store)}
}

func (r TicketAggregateRootRepository) Load(ctx context.Context, aggregateID string) (*es.AggregateRoot, error) {
	root, err := r.AggregateRepository.Load(ctx, aggregateID)
	if err != nil {
		if err == es.ErrAggregateNotFound {
			return nil, errors.Wrap(errors.ErrNotFound, "ticket not found")
		}
		return nil, err
	}

	return root, nil
}
