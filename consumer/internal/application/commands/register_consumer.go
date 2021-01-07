package commands

import (
	"context"

	"github.com/stackus/ftgogo/consumer/internal/domain"
)

type RegisterConsumer struct {
	Name string
}

type RegisterConsumerHandler struct {
	repo      domain.ConsumerRepository
	publisher domain.ConsumerPublisher
}

func NewRegisterConsumerHandler(consumerRepo domain.ConsumerRepository, consumerPublisher domain.ConsumerPublisher) RegisterConsumerHandler {
	return RegisterConsumerHandler{
		repo:      consumerRepo,
		publisher: consumerPublisher,
	}
}

func (h RegisterConsumerHandler) Handle(ctx context.Context, cmd RegisterConsumer) (string, error) {
	root, err := h.repo.Save(ctx, &domain.RegisterConsumer{
		Name: cmd.Name,
	})
	if err != nil {
		return "", err
	}

	return root.AggregateID(), h.publisher.PublishEntityEvents(ctx, root)
}
