package ports

import (
	"context"

	"github.com/stackus/edat/core"
	"github.com/stackus/edat/es"

	"github.com/stackus/ftgogo/accounting/internal/domain"
)

type AccountRepository interface {
	Load(ctx context.Context, aggregateID string) (*domain.Account, error)
	Save(ctx context.Context, command core.Command, options ...es.AggregateRootOption) (*domain.Account, error)
	Update(ctx context.Context, aggregateID string, command core.Command, options ...es.AggregateRootOption) (*domain.Account, error)
}
