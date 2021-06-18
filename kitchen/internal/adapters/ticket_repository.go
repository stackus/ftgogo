package adapters

import (
	"context"

	"github.com/stackus/edat/core"
	"github.com/stackus/edat/es"

	"github.com/stackus/ftgogo/kitchen/internal/domain"
)

type TicketRepository interface {
	Load(ctx context.Context, aggregateID string) (*domain.Ticket, error)
	Save(ctx context.Context, command core.Command, options ...es.AggregateRootOption) (*domain.Ticket, error)
	Update(ctx context.Context, aggregateID string, command core.Command, options ...es.AggregateRootOption) (*domain.Ticket, error)
}
