package adapters

import (
	"context"

	"github.com/stackus/edat/core"
	"github.com/stackus/edat/es"

	"github.com/stackus/ftgogo/accounting/internal/application/ports"
	"github.com/stackus/ftgogo/accounting/internal/domain"
)

type AccountAggregateRepository struct {
	store es.AggregateRepository
}

var _ ports.AccountRepository = (*AccountAggregateRepository)(nil)

func NewAccountAggregateRepository(store es.AggregateRootStore) *AccountAggregateRepository {
	return &AccountAggregateRepository{store: es.NewAggregateRootRepository(domain.NewAccount, store)}
}

func (r AccountAggregateRepository) Load(ctx context.Context, aggregateID string) (*domain.Account, error) {
	root, err := r.store.Load(ctx, aggregateID)
	if err != nil {
		if err == es.ErrAggregateNotFound {
			return nil, domain.ErrAccountNotFound
		}
		return nil, err
	}

	return root.Aggregate().(*domain.Account), nil
}

func (r AccountAggregateRepository) Save(ctx context.Context, command core.Command, options ...es.AggregateRootOption) (*domain.Account, error) {
	root, err := r.store.Save(ctx, command, options...)
	if err != nil {
		return nil, err
	}
	return root.Aggregate().(*domain.Account), nil
}

func (r AccountAggregateRepository) Update(ctx context.Context, aggregateID string, command core.Command, options ...es.AggregateRootOption) (*domain.Account, error) {
	root, err := r.store.Update(ctx, aggregateID, command, options...)
	if err != nil {
		if err == es.ErrAggregateNotFound {
			return nil, domain.ErrAccountNotFound
		}
		return nil, err
	}
	return root.Aggregate().(*domain.Account), nil
}
