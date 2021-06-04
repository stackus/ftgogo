package commands

import (
	"context"

	"github.com/stackus/ftgogo/accounting/internal/adapters"
	"github.com/stackus/ftgogo/accounting/internal/domain"
)

type EnableAccount struct {
	AccountID string
}

type EnableAccountHandler struct {
	repo adapters.AccountRepository
}

func NewEnableAccountHandler(accountRepo adapters.AccountRepository) EnableAccountHandler {
	return EnableAccountHandler{repo: accountRepo}
}

func (h EnableAccountHandler) Handle(ctx context.Context, cmd EnableAccount) error {
	_, err := h.repo.Update(ctx, cmd.AccountID, &domain.EnableAccount{})
	return err
}
