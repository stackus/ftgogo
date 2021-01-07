package commands

import (
	"context"

	"github.com/stackus/ftgogo/account/internal/domain"
)

type EnableAccount struct {
	AccountID string
}

type EnableAccountHandler struct {
	repo domain.AccountRepository
}

func NewEnableAccountHandler(accountRepo domain.AccountRepository) EnableAccountHandler {
	return EnableAccountHandler{repo: accountRepo}
}

func (h EnableAccountHandler) Handle(ctx context.Context, cmd EnableAccount) error {
	_, err := h.repo.Update(ctx, cmd.AccountID, &domain.EnableAccount{})
	return err
}
