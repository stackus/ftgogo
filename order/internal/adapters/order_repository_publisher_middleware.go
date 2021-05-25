package adapters

import (
	"context"

	"github.com/stackus/edat/core"
	"github.com/stackus/edat/es"

	"github.com/stackus/ftgogo/order/internal/domain"
)

type orderRepositoryPublisherMiddleware struct {
	domain.OrderRepository
	publisher domain.OrderPublisher
}

var _ domain.OrderRepository = (*orderRepositoryPublisherMiddleware)(nil)

func NewOrderRepositoryPublisherMiddleware(repository domain.OrderRepository, publisher domain.OrderPublisher) domain.OrderRepository {
	return &orderRepositoryPublisherMiddleware{
		OrderRepository: repository,
		publisher:       publisher,
	}
}

func (r orderRepositoryPublisherMiddleware) Save(ctx context.Context, command core.Command, options ...es.AggregateRootOption) (*es.AggregateRoot, error) {
	root, err := r.OrderRepository.Save(ctx, command, options...)
	if err != nil {
		return root, err
	}

	return root, r.publisher.PublishEntityEvents(ctx, root)
}

func (r orderRepositoryPublisherMiddleware) Update(ctx context.Context, aggregateID string, command core.Command, options ...es.AggregateRootOption) (*es.AggregateRoot, error) {
	root, err := r.OrderRepository.Update(ctx, aggregateID, command, options...)
	if err != nil {
		return root, err
	}

	return root, r.publisher.PublishEntityEvents(ctx, root)
}
