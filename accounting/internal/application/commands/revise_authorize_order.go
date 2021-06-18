package commands

import (
	"context"

	"github.com/stackus/ftgogo/accounting/internal/application/ports"
	"github.com/stackus/ftgogo/accounting/internal/domain"
)

type ReviseAuthorizeOrder struct {
	ConsumerID string
	OrderID    string
	OrderTotal int
}

type ReviseAuthorizeOrderHandler struct {
	repo ports.AccountRepository
}

func NewReviseAuthorizeOrderHandler(accountRepo ports.AccountRepository) ReviseAuthorizeOrderHandler {
	return ReviseAuthorizeOrderHandler{repo: accountRepo}
}

func (h ReviseAuthorizeOrderHandler) Handle(ctx context.Context, cmd ReviseAuthorizeOrder) error {
	_, err := h.repo.Update(ctx, cmd.ConsumerID, &domain.ReviseAuthorizeOrder{
		OrderID:    cmd.OrderID,
		OrderTotal: cmd.OrderTotal,
	})

	return err
}
