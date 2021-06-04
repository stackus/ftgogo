package commands

import (
	"context"

	"github.com/stackus/ftgogo/consumer/internal/adapters"
	"github.com/stackus/ftgogo/consumer/internal/domain"
)

type UpdateConsumer struct {
	ConsumerID string
	Name       string
}

type UpdateConsumerHandler struct {
	repo adapters.ConsumerRepository
}

func NewUpdateConsumerHandler(repo adapters.ConsumerRepository) UpdateConsumerHandler {
	return UpdateConsumerHandler{repo: repo}
}

func (h UpdateConsumerHandler) Handle(ctx context.Context, cmd UpdateConsumer) error {
	_, err := h.repo.Update(ctx, cmd.ConsumerID, &domain.UpdateConsumer{Name: cmd.Name})
	return err
}
