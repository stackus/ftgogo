package queries

import (
	"context"

	"github.com/stackus/ftgogo/web-bff/internal/domain"
)

type GetConsumer struct {
	ConsumerID string
}

type GetConsumerHandler struct {
	repo domain.ConsumerRepository
}

func NewGetConsumerHandler(repo domain.ConsumerRepository) GetConsumerHandler {
	return GetConsumerHandler{repo: repo}
}

func (h GetConsumerHandler) Handle(ctx context.Context, cmd GetConsumer) (*domain.Consumer, error) {
	return h.repo.Find(ctx, domain.FindConsumer{
		ConsumerID: cmd.ConsumerID,
	})
}
