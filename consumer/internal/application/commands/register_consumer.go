package commands

import (
	"context"

	"github.com/stackus/ftgogo/consumer/internal/adapters"
	"github.com/stackus/ftgogo/consumer/internal/domain"
)

type RegisterConsumer struct {
	Name string
}

type RegisterConsumerHandler struct {
	repo adapters.ConsumerRepository
}

func NewRegisterConsumerHandler(repo adapters.ConsumerRepository) RegisterConsumerHandler {
	return RegisterConsumerHandler{
		repo: repo,
	}
}

func (h RegisterConsumerHandler) Handle(ctx context.Context, cmd RegisterConsumer) (string, error) {
	root, err := h.repo.Save(ctx, &domain.RegisterConsumer{
		Name: cmd.Name,
	})
	if err != nil {
		return "", err
	}

	return root.AggregateID(), nil
}
