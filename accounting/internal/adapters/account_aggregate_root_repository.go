package adapters

import (
	"context"

	"github.com/stackus/edat/es"
	"github.com/stackus/errors"

	"github.com/stackus/ftgogo/accounting/internal/domain"
)

type AccountAggregateRootRepository struct {
	es.AggregateRepository
}

func NewAccountAggregateRootRepository(store es.AggregateRootStore) *AccountAggregateRootRepository {
	return &AccountAggregateRootRepository{es.NewAggregateRootRepository(domain.NewAccount, store)}
}

func (r AccountAggregateRootRepository) Load(ctx context.Context, aggregateID string) (*es.AggregateRoot, error) {
	root, err := r.AggregateRepository.Load(ctx, aggregateID)
	if err != nil {
		if err == es.ErrAggregateNotFound {
			return nil, errors.Wrap(errors.ErrNotFound, "account not found")
		}
		return nil, err
	}

	return root, nil
}
