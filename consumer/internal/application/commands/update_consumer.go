package commands

import (
	"context"

	"github.com/stackus/ftgogo/consumer/internal/application/ports"
	"github.com/stackus/ftgogo/consumer/internal/domain"
)

type UpdateConsumer struct {
	ConsumerID string
	Name       string
}

type UpdateConsumerHandler struct {
	repo ports.ConsumerRepository
}

func NewUpdateConsumerHandler(repo ports.ConsumerRepository) UpdateConsumerHandler {
	return UpdateConsumerHandler{repo: repo}
}

func (h UpdateConsumerHandler) Handle(ctx context.Context, cmd UpdateConsumer) error {
	_, err := h.repo.Update(ctx, cmd.ConsumerID, &domain.UpdateConsumer{Name: cmd.Name})
	return err
}
