package commands

import (
	"context"

	"github.com/stackus/ftgogo/store-web/internal/adapters"
)

type DisableAccount struct {
	AccountID string
}

type DisableAccountHandler struct {
	repo adapters.AccountingRepository
}

func NewDisableAccountHandler(repo adapters.AccountingRepository) DisableAccountHandler {
	return DisableAccountHandler{repo: repo}
}

func (h DisableAccountHandler) Handle(ctx context.Context, cmd DisableAccount) error {
	return h.repo.Disable(ctx, adapters.DisableAccount{AccountID: cmd.AccountID})
}
