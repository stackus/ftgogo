package adapters

import (
	"context"

	"github.com/stackus/edat/es"
	"github.com/stackus/errors"

	"github.com/stackus/ftgogo/order/internal/domain"
)

type OrderAggregateRootRepository struct {
	es.AggregateRepository
}

func NewOrderAggregateRootRepository(store es.AggregateRootStore) *OrderAggregateRootRepository {
	return &OrderAggregateRootRepository{es.NewAggregateRootRepository(domain.NewOrder, store)}
}

func (r OrderAggregateRootRepository) Load(ctx context.Context, aggregateID string) (*es.AggregateRoot, error) {
	root, err := r.AggregateRepository.Load(ctx, aggregateID)
	if err != nil {
		if err == es.ErrAggregateNotFound {
			return nil, errors.Wrap(errors.ErrNotFound, "order not found")
		}
		return nil, err
	}

	return root, nil
}
