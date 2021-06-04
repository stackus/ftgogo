package commands

import (
	"context"

	"github.com/stackus/ftgogo/accounting/internal/adapters"
	"github.com/stackus/ftgogo/accounting/internal/domain"
)

type DisableAccount struct {
	AccountID string
}

type DisableAccountHandler struct {
	repo adapters.AccountRepository
}

func NewDisableAccountHandler(accountRepo adapters.AccountRepository) DisableAccountHandler {
	return DisableAccountHandler{repo: accountRepo}
}

func (h DisableAccountHandler) Handle(ctx context.Context, cmd DisableAccount) error {
	_, err := h.repo.Update(ctx, cmd.AccountID, &domain.DisableAccount{})
	return err
}
