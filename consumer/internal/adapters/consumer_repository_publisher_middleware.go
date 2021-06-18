package adapters

import (
	"context"

	"github.com/stackus/edat/core"
	"github.com/stackus/edat/es"

	"github.com/stackus/ftgogo/consumer/internal/domain"
)

type consumerRepositoryPublisherMiddleware struct {
	ConsumerRepository
	publisher ConsumerPublisher
}

var _ ConsumerRepository = (*consumerRepositoryPublisherMiddleware)(nil)

func NewConsumerRepositoryPublisherMiddleware(repository ConsumerRepository, publisher ConsumerPublisher) ConsumerRepository {
	return &consumerRepositoryPublisherMiddleware{
		ConsumerRepository: repository,
		publisher:          publisher,
	}
}

func (r consumerRepositoryPublisherMiddleware) Save(ctx context.Context, command core.Command, options ...es.AggregateRootOption) (*domain.Consumer, error) {
	consumer, err := r.ConsumerRepository.Save(ctx, command, options...)
	if err != nil {
		return consumer, err
	}

	return consumer, r.publisher.PublishEntityEvents(ctx, consumer)
}

func (r consumerRepositoryPublisherMiddleware) Update(ctx context.Context, aggregateID string, command core.Command, options ...es.AggregateRootOption) (*domain.Consumer, error) {
	consumer, err := r.ConsumerRepository.Update(ctx, aggregateID, command, options...)
	if err != nil {
		return consumer, err
	}

	return consumer, r.publisher.PublishEntityEvents(ctx, consumer)
}
