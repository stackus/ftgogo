package commands

import (
	"context"

	"github.com/stackus/ftgogo/store-web/internal/adapters"
)

type EnableAccount struct {
	AccountID string
}

type EnableAccountHandler struct {
	repo adapters.AccountingRepository
}

func NewEnableAccountHandler(repo adapters.AccountingRepository) EnableAccountHandler {
	return EnableAccountHandler{repo: repo}
}

func (h EnableAccountHandler) Handle(ctx context.Context, cmd EnableAccount) error {
	return h.repo.Enable(ctx, adapters.EnableAccount{AccountID: cmd.AccountID})
}
