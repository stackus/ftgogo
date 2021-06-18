package adapters

import (
	"context"

	"github.com/stackus/edat/core"
	"github.com/stackus/edat/es"

	"github.com/stackus/ftgogo/order/internal/domain"
)

type orderRepositoryPublisherMiddleware struct {
	OrderRepository
	publisher OrderPublisher
}

var _ OrderRepository = (*orderRepositoryPublisherMiddleware)(nil)

func NewOrderRepositoryPublisherMiddleware(repository OrderRepository, publisher OrderPublisher) OrderRepository {
	return &orderRepositoryPublisherMiddleware{
		OrderRepository: repository,
		publisher:       publisher,
	}
}

func (r orderRepositoryPublisherMiddleware) Save(ctx context.Context, command core.Command, options ...es.AggregateRootOption) (*domain.Order, error) {
	order, err := r.OrderRepository.Save(ctx, command, options...)
	if err != nil {
		return order, err
	}

	return order, r.publisher.PublishEntityEvents(ctx, order)
}

func (r orderRepositoryPublisherMiddleware) Update(ctx context.Context, aggregateID string, command core.Command, options ...es.AggregateRootOption) (*domain.Order, error) {
	order, err := r.OrderRepository.Update(ctx, aggregateID, command, options...)
	if err != nil {
		return order, err
	}

	return order, r.publisher.PublishEntityEvents(ctx, order)
}
