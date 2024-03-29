package commands

import (
	"context"

	"github.com/stackus/ftgogo/consumer/internal/application/ports"
	"github.com/stackus/ftgogo/consumer/internal/domain"
)

type ValidateOrderByConsumer struct {
	ConsumerID string
	OrderID    string
	OrderTotal int
}

type ValidateOrderByConsumerHandler struct {
	repo ports.ConsumerRepository
}

func NewValidateOrderByConsumerHandler(consumerRepo ports.ConsumerRepository) ValidateOrderByConsumerHandler {
	return ValidateOrderByConsumerHandler{repo: consumerRepo}
}

func (h ValidateOrderByConsumerHandler) Handle(ctx context.Context, cmd ValidateOrderByConsumer) error {
	consumer, err := h.repo.Load(ctx, cmd.ConsumerID)
	if err != nil {
		return domain.ErrConsumerNotFound
	}

	err = consumer.ValidateOrderByConsumer(cmd.OrderTotal)
	if err != nil {
		return domain.ErrOrderNotValidated
	}

	return nil
}
