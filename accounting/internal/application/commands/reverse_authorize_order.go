package commands

import (
	"context"

	"github.com/stackus/ftgogo/accounting/internal/application/ports"
	"github.com/stackus/ftgogo/accounting/internal/domain"
)

type ReverseAuthorizeOrder struct {
	ConsumerID string
	OrderID    string
	OrderTotal int
}

type ReverseAuthorizeOrderHandler struct {
	repo ports.AccountRepository
}

func NewReverseAuthorizeOrderHandler(accountRepo ports.AccountRepository) ReverseAuthorizeOrderHandler {
	return ReverseAuthorizeOrderHandler{repo: accountRepo}
}

func (h ReverseAuthorizeOrderHandler) Handle(ctx context.Context, cmd ReverseAuthorizeOrder) error {
	_, err := h.repo.Update(ctx, cmd.ConsumerID, &domain.ReverseAuthorizeOrder{
		OrderID:    cmd.OrderID,
		OrderTotal: cmd.OrderTotal,
	})

	return err
}
