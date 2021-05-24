package commands

import (
	"context"

	"github.com/stackus/ftgogo/web-bff/internal/domain"
)

type RegisterConsumer struct {
	Name string
}

type RegisterConsumerHandler struct {
	repo domain.ConsumerRepository
}

func NewRegisterConsumerHandler(repo domain.ConsumerRepository) RegisterConsumerHandler {
	return RegisterConsumerHandler{repo: repo}
}

func (h RegisterConsumerHandler) Handle(ctx context.Context, cmd RegisterConsumer) (string, error) {
	consumerID, err := h.repo.Register(ctx, cmd.Name)
	if err != nil {
		return "", err
	}

	return consumerID, nil
}
