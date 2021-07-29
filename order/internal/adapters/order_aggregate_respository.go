package adapters

import (
	"context"

	"github.com/stackus/edat/core"
	"github.com/stackus/edat/es"
	"github.com/stackus/errors"

	"github.com/stackus/ftgogo/order/internal/application/ports"
	"github.com/stackus/ftgogo/order/internal/domain"
)

type OrderAggregateRepository struct {
	store es.AggregateRepository
}

var _ ports.OrderRepository = (*OrderAggregateRepository)(nil)

func NewOrderAggregateRepository(store es.AggregateRootStore) *OrderAggregateRepository {
	return &OrderAggregateRepository{store: es.NewAggregateRootRepository(domain.NewOrder, store)}
}

func (r OrderAggregateRepository) Load(ctx context.Context, aggregateID string) (*domain.Order, error) {
	root, err := r.store.Load(ctx, aggregateID)
	if err != nil {
		if err == es.ErrAggregateNotFound {
			return nil, errors.Wrap(errors.ErrNotFound, "order not found")
		}
		return nil, err
	}

	return root.Aggregate().(*domain.Order), nil
}

func (r OrderAggregateRepository) Save(ctx context.Context, command core.Command, options ...es.AggregateRootOption) (*domain.Order, error) {
	root, err := r.store.Save(ctx, command, options...)
	if err != nil {
		return nil, err
	}
	return root.Aggregate().(*domain.Order), nil
}

func (r OrderAggregateRepository) Update(ctx context.Context, aggregateID string, command core.Command, options ...es.AggregateRootOption) (*domain.Order, error) {
	root, err := r.store.Update(ctx, aggregateID, command, options...)
	if err != nil {
		return nil, err
	}
	return root.Aggregate().(*domain.Order), nil
}
