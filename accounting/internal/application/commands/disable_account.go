package commands

import (
	"context"

	"github.com/stackus/ftgogo/accounting/internal/application/ports"
	"github.com/stackus/ftgogo/accounting/internal/domain"
)

type DisableAccount struct {
	AccountID string
}

type DisableAccountHandler struct {
	repo ports.AccountRepository
}

func NewDisableAccountHandler(accountRepo ports.AccountRepository) DisableAccountHandler {
	return DisableAccountHandler{repo: accountRepo}
}

func (h DisableAccountHandler) Handle(ctx context.Context, cmd DisableAccount) error {
	_, err := h.repo.Update(ctx, cmd.AccountID, &domain.DisableAccount{})
	return err
}
