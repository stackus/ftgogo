package queries

import (
	"context"

	"github.com/stackus/ftgogo/consumer/internal/application/ports"
	"github.com/stackus/ftgogo/consumer/internal/domain"
)

type GetConsumer struct {
	ConsumerID string
}

type GetConsumerHandler struct {
	repo ports.ConsumerRepository
}

func NewGetConsumerHandler(consumerRepo ports.ConsumerRepository) GetConsumerHandler {
	return GetConsumerHandler{repo: consumerRepo}
}

func (h GetConsumerHandler) Handle(ctx context.Context, query GetConsumer) (*domain.Consumer, error) {
	return h.repo.Load(ctx, query.ConsumerID)
}
