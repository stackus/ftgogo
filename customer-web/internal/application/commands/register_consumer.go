package commands

import (
	"context"

	"github.com/stackus/ftgogo/customer-web/internal/adapters"
)

type RegisterConsumer struct {
	Name string
}

type RegisterConsumerHandler struct {
	repo adapters.ConsumerRepository
}

func NewRegisterConsumerHandler(repo adapters.ConsumerRepository) RegisterConsumerHandler {
	return RegisterConsumerHandler{repo: repo}
}

func (h RegisterConsumerHandler) Handle(ctx context.Context, cmd RegisterConsumer) (string, error) {
	consumerID, err := h.repo.Register(ctx, adapters.RegisterConsumer{
		Name: cmd.Name,
	})
	if err != nil {
		return "", err
	}

	return consumerID, nil
}
