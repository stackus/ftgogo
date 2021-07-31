package queries

import (
	"context"

	"github.com/stackus/ftgogo/store-web/internal/adapters"
	"github.com/stackus/ftgogo/store-web/internal/domain"
)

type GetAccount struct {
	AccountID string
}

type GetAccountHandler struct {
	repo adapters.AccountingRepository
}

func NewGetAccountHandler(repo adapters.AccountingRepository) GetAccountHandler {
	return GetAccountHandler{repo: repo}
}

func (h GetAccountHandler) Handle(ctx context.Context, cmd GetAccount) (*domain.Account, error) {
	return h.repo.Find(ctx, adapters.FindAccount{
		AccountID: cmd.AccountID,
	})
}
