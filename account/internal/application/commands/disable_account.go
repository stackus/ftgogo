package commands

import (
	"context"

	"github.com/stackus/ftgogo/account/internal/domain"
)

type DisableAccount struct {
	AccountID string
}

type DisableAccountHandler struct {
	repo domain.AccountRepository
}

func NewDisableAccountHandler(accountRepo domain.AccountRepository) DisableAccountHandler {
	return DisableAccountHandler{repo: accountRepo}
}

func (h DisableAccountHandler) Handle(ctx context.Context, cmd DisableAccount) error {
	_, err := h.repo.Update(ctx, cmd.AccountID, &domain.DisableAccount{})
	return err
}
