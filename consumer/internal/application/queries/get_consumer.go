package queries

import (
	"context"

	"github.com/stackus/ftgogo/consumer/internal/adapters"
	"github.com/stackus/ftgogo/consumer/internal/domain"
)

type GetConsumer struct {
	ConsumerID string
}

type GetConsumerHandler struct {
	repo adapters.ConsumerRepository
}

func NewGetConsumerHandler(consumerRepo adapters.ConsumerRepository) GetConsumerHandler {
	return GetConsumerHandler{repo: consumerRepo}
}

func (h GetConsumerHandler) Handle(ctx context.Context, query GetConsumer) (*domain.Consumer, error) {
	root, err := h.repo.Load(ctx, query.ConsumerID)
	if err != nil {
		return nil, err
	}

	return root.Aggregate().(*domain.Consumer), nil
}
