package commands

import (
	"context"

	"github.com/stackus/ftgogo/account/internal/domain"
)

type AuthorizeOrder struct {
	ConsumerID string
	OrderID    string
	OrderTotal int
}

type AuthorizeOrderHandler struct {
	repo domain.AccountRepository
}

func NewAuthorizeOrderHandler(accountRepo domain.AccountRepository) AuthorizeOrderHandler {
	return AuthorizeOrderHandler{repo: accountRepo}
}

func (h AuthorizeOrderHandler) Handle(ctx context.Context, cmd AuthorizeOrder) error {
	_, err := h.repo.Update(ctx, cmd.ConsumerID, &domain.AuthorizeOrder{
		OrderID:    cmd.OrderID,
		OrderTotal: cmd.OrderTotal,
	})

	return err
}
