package adapters

import (
	"context"

	"github.com/stackus/edat/core"
	"github.com/stackus/edat/es"

	"github.com/stackus/ftgogo/order/internal/domain"
)

type OrderRepository interface {
	Load(ctx context.Context, aggregateID string) (*domain.Order, error)
	Save(ctx context.Context, command core.Command, options ...es.AggregateRootOption) (*domain.Order, error)
	Update(ctx context.Context, aggregateID string, command core.Command, options ...es.AggregateRootOption) (*domain.Order, error)
}
