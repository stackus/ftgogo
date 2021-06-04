package commands

import (
	"context"

	"github.com/stackus/ftgogo/consumer/internal/adapters"
	"github.com/stackus/ftgogo/consumer/internal/domain"
)

type ValidateOrderByConsumer struct {
	ConsumerID string
	OrderID    string
	OrderTotal int
}

type ValidateOrderByConsumerHandler struct {
	repo adapters.ConsumerRepository
}

func NewValidateOrderByConsumerHandler(consumerRepo adapters.ConsumerRepository) ValidateOrderByConsumerHandler {
	return ValidateOrderByConsumerHandler{repo: consumerRepo}
}

func (h ValidateOrderByConsumerHandler) Handle(ctx context.Context, cmd ValidateOrderByConsumer) error {
	root, err := h.repo.Load(ctx, cmd.ConsumerID)
	if err != nil {
		return domain.ErrConsumerNotFound
	}

	consumer := root.Aggregate().(*domain.Consumer)

	err = consumer.ValidateOrderByConsumer(cmd.OrderTotal)
	if err != nil {
		return domain.ErrOrderNotValidated
	}

	return nil
}
