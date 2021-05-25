package adapters

import (
	"context"

	"github.com/stackus/edat/core"
	"github.com/stackus/edat/es"

	"github.com/stackus/ftgogo/kitchen/internal/domain"
)

type ticketRepositoryPublisherMiddleware struct {
	domain.TicketRepository
	publisher domain.TicketPublisher
}

var _ domain.TicketRepository = (*ticketRepositoryPublisherMiddleware)(nil)

func NewTicketRepositoryPublisherMiddleware(repository domain.TicketRepository, publisher domain.TicketPublisher) domain.TicketRepository {
	return &ticketRepositoryPublisherMiddleware{
		TicketRepository: repository,
		publisher:        publisher,
	}
}

func (r ticketRepositoryPublisherMiddleware) Save(ctx context.Context, command core.Command, options ...es.AggregateRootOption) (*es.AggregateRoot, error) {
	root, err := r.TicketRepository.Save(ctx, command, options...)
	if err != nil {
		return root, err
	}

	return root, r.publisher.PublishEntityEvents(ctx, root)
}

func (r ticketRepositoryPublisherMiddleware) Update(ctx context.Context, aggregateID string, command core.Command, options ...es.AggregateRootOption) (*es.AggregateRoot, error) {
	root, err := r.TicketRepository.Update(ctx, aggregateID, command, options...)
	if err != nil {
		return root, err
	}

	return root, r.publisher.PublishEntityEvents(ctx, root)
}
