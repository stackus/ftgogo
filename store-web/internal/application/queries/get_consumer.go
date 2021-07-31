package queries

import (
	"context"

	"github.com/stackus/ftgogo/store-web/internal/adapters"
	"github.com/stackus/ftgogo/store-web/internal/domain"
)

type GetConsumer struct {
	ConsumerID string
}

type GetConsumerHandler struct {
	repo adapters.ConsumerRepository
}

func NewGetConsumerHandler(repo adapters.ConsumerRepository) GetConsumerHandler {
	return GetConsumerHandler{repo: repo}
}

func (h GetConsumerHandler) Handle(ctx context.Context, cmd GetConsumer) (*domain.Consumer, error) {
	return h.repo.Find(ctx, adapters.FindConsumer{
		ConsumerID: cmd.ConsumerID,
	})
}
