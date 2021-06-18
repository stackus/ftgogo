package adapters

import (
	"context"

	"github.com/stackus/edat/core"
	"github.com/stackus/edat/es"
	"github.com/stackus/errors"

	"github.com/stackus/ftgogo/consumer/internal/domain"
)

type ConsumerAggregateRepository struct {
	store es.AggregateRepository
}

var _ ConsumerRepository = (*ConsumerAggregateRepository)(nil)

func NewConsumerAggregateRepository(store es.AggregateRootStore) *ConsumerAggregateRepository {
	return &ConsumerAggregateRepository{store: es.NewAggregateRootRepository(domain.NewConsumer, store)}
}

func (r ConsumerAggregateRepository) Load(ctx context.Context, aggregateID string) (*domain.Consumer, error) {
	root, err := r.store.Load(ctx, aggregateID)
	if err != nil {
		if err == es.ErrAggregateNotFound {
			return nil, errors.Wrap(errors.ErrNotFound, "consumer not found")
		}
		return nil, err
	}

	return root.Aggregate().(*domain.Consumer), nil
}

func (r ConsumerAggregateRepository) Save(ctx context.Context, command core.Command, options ...es.AggregateRootOption) (*domain.Consumer, error) {
	root, err := r.store.Save(ctx, command, options...)
	if err != nil {
		return nil, err
	}
	return root.Aggregate().(*domain.Consumer), nil
}

func (r ConsumerAggregateRepository) Update(ctx context.Context, aggregateID string, command core.Command, options ...es.AggregateRootOption) (*domain.Consumer, error) {
	root, err := r.store.Update(ctx, aggregateID, command, options...)
	if err != nil {
		return nil, err
	}
	return root.Aggregate().(*domain.Consumer), nil
}
