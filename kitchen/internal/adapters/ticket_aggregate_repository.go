package adapters

import (
	"context"

	"github.com/stackus/edat/core"
	"github.com/stackus/edat/es"
	"github.com/stackus/errors"

	"github.com/stackus/ftgogo/kitchen/internal/domain"
)

type TicketAggregateRepository struct {
	store es.AggregateRepository
}

var _ TicketRepository = (*TicketAggregateRepository)(nil)

func NewTicketAggregateRepository(store es.AggregateRootStore) *TicketAggregateRepository {
	return &TicketAggregateRepository{store: es.NewAggregateRootRepository(domain.NewTicket, store)}
}

func (r TicketAggregateRepository) Load(ctx context.Context, aggregateID string) (*domain.Ticket, error) {
	root, err := r.store.Load(ctx, aggregateID)
	if err != nil {
		if err == es.ErrAggregateNotFound {
			return nil, errors.Wrap(errors.ErrNotFound, "ticket not found")
		}
		return nil, err
	}

	return root.Aggregate().(*domain.Ticket), nil
}

func (r TicketAggregateRepository) Save(ctx context.Context, command core.Command, options ...es.AggregateRootOption) (*domain.Ticket, error) {
	root, err := r.store.Save(ctx, command, options...)
	if err != nil {
		return nil, err
	}
	return root.Aggregate().(*domain.Ticket), nil
}

func (r TicketAggregateRepository) Update(ctx context.Context, aggregateID string, command core.Command, options ...es.AggregateRootOption) (*domain.Ticket, error) {
	root, err := r.store.Update(ctx, aggregateID, command, options...)
	if err != nil {
		return nil, err
	}
	return root.Aggregate().(*domain.Ticket), nil
}
