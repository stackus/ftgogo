package commands

import (
	"context"

	"github.com/stackus/ftgogo/account/internal/domain"
)

type ReviseAuthorizeOrder struct {
	ConsumerID string
	OrderID    string
	OrderTotal int
}

type ReviseAuthorizeOrderHandler struct {
	repo domain.AccountRepository
}

func NewReviseAuthorizeOrderHandler(accountRepo domain.AccountRepository) ReviseAuthorizeOrderHandler {
	return ReviseAuthorizeOrderHandler{repo: accountRepo}
}

func (h ReviseAuthorizeOrderHandler) Handle(ctx context.Context, cmd ReviseAuthorizeOrder) error {
	_, err := h.repo.Update(ctx, cmd.ConsumerID, &domain.ReviseAuthorizeOrder{
		OrderID:    cmd.OrderID,
		OrderTotal: cmd.OrderTotal,
	})

	return err
}
