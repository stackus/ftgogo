package adapters

import (
	"context"

	"github.com/stackus/edat/core"
	"github.com/stackus/edat/es"

	"github.com/stackus/ftgogo/kitchen/internal/domain"
)

type ticketRepositoryPublisherMiddleware struct {
	TicketRepository
	publisher TicketPublisher
}

var _ TicketRepository = (*ticketRepositoryPublisherMiddleware)(nil)

func NewTicketRepositoryPublisherMiddleware(repository TicketRepository, publisher TicketPublisher) TicketRepository {
	return &ticketRepositoryPublisherMiddleware{
		TicketRepository: repository,
		publisher:        publisher,
	}
}

func (r ticketRepositoryPublisherMiddleware) Save(ctx context.Context, command core.Command, options ...es.AggregateRootOption) (*domain.Ticket, error) {
	ticket, err := r.TicketRepository.Save(ctx, command, options...)
	if err != nil {
		return ticket, err
	}

	return ticket, r.publisher.PublishEntityEvents(ctx, ticket)
}

func (r ticketRepositoryPublisherMiddleware) Update(ctx context.Context, aggregateID string, command core.Command, options ...es.AggregateRootOption) (*domain.Ticket, error) {
	ticket, err := r.TicketRepository.Update(ctx, aggregateID, command, options...)
	if err != nil {
		return ticket, err
	}

	return ticket, r.publisher.PublishEntityEvents(ctx, ticket)
}
