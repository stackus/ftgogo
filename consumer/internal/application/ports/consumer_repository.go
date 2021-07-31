package ports

import (
	"context"

	"github.com/stackus/edat/core"
	"github.com/stackus/edat/es"

	"github.com/stackus/ftgogo/consumer/internal/domain"
)

type ConsumerRepository interface {
	Load(ctx context.Context, aggregateID string) (*domain.Consumer, error)
	Save(ctx context.Context, command core.Command, options ...es.AggregateRootOption) (*domain.Consumer, error)
	Update(ctx context.Context, aggregateID string, command core.Command, options ...es.AggregateRootOption) (*domain.Consumer, error)
}
