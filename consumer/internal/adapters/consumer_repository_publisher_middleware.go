package adapters

import (
	"context"
	"github.com/stackus/edat/core"
	"github.com/stackus/edat/es"
	"github.com/stackus/ftgogo/consumer/internal/domain"
)

type consumerRepositoryPublisherMiddleware struct {
	domain.ConsumerRepository
	publisher domain.ConsumerPublisher
}

var _ domain.ConsumerRepository = (*consumerRepositoryPublisherMiddleware)(nil)

func NewConsumerRepositoryPublisherMiddleware(repository domain.ConsumerRepository, publisher domain.ConsumerPublisher) domain.ConsumerRepository {
	return &consumerRepositoryPublisherMiddleware{
		ConsumerRepository: repository,
		publisher:          publisher,
	}
}

func (r consumerRepositoryPublisherMiddleware) Save(ctx context.Context, command core.Command, options ...es.AggregateRootOption) (*es.AggregateRoot, error) {
	root, err := r.ConsumerRepository.Save(ctx, command, options...)
	if err != nil {
		return root, err
	}

	return root, r.publisher.PublishEntityEvents(ctx, root)
}

func (r consumerRepositoryPublisherMiddleware) Update(ctx context.Context, aggregateID string, command core.Command, options ...es.AggregateRootOption) (*es.AggregateRoot, error) {
	root, err := r.ConsumerRepository.Update(ctx, aggregateID, command, options...)
	if err != nil {
		return root, err
	}

	return root, r.publisher.PublishEntityEvents(ctx, root)
}
