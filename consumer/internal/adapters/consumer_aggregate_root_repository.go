package adapters

import (
	"context"

	"github.com/stackus/edat/es"
	"github.com/stackus/errors"

	"github.com/stackus/ftgogo/consumer/internal/domain"
)

type ConsumerAggregateRootRepository struct {
	es.AggregateRepository
}

func NewConsumerAggregateRootRepository(store es.AggregateRootStore) *ConsumerAggregateRootRepository {
	return &ConsumerAggregateRootRepository{es.NewAggregateRootRepository(domain.NewConsumer, store)}
}

func (r ConsumerAggregateRootRepository) Load(ctx context.Context, aggregateID string) (*es.AggregateRoot, error) {
	root, err := r.AggregateRepository.Load(ctx, aggregateID)
	if err != nil {
		if err == es.ErrAggregateNotFound {
			return nil, errors.Wrap(errors.ErrNotFound, "consumer not found")
		}
		return nil, err
	}

	return root, nil
}
