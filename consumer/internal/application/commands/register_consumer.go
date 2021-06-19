package commands

import (
	"context"

	"github.com/stackus/ftgogo/consumer/internal/application/ports"
	"github.com/stackus/ftgogo/consumer/internal/domain"
)

type RegisterConsumer struct {
	Name string
}

type RegisterConsumerHandler struct {
	repo ports.ConsumerRepository
}

func NewRegisterConsumerHandler(repo ports.ConsumerRepository) RegisterConsumerHandler {
	return RegisterConsumerHandler{
		repo: repo,
	}
}

func (h RegisterConsumerHandler) Handle(ctx context.Context, cmd RegisterConsumer) (string, error) {
	consumer, err := h.repo.Save(ctx, &domain.RegisterConsumer{
		Name: cmd.Name,
	})
	if err != nil {
		return "", err
	}

	return consumer.ID(), nil
}
