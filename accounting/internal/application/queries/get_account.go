package queries

import (
	"context"

	"github.com/stackus/ftgogo/accounting/internal/application/ports"
	"github.com/stackus/ftgogo/accounting/internal/domain"
)

type GetAccount struct {
	AccountID string
}

type GetAccountHandler struct {
	repo ports.AccountRepository
}

func NewGetAccountHandler(accountRepo ports.AccountRepository) GetAccountHandler {
	return GetAccountHandler{repo: accountRepo}
}

func (h GetAccountHandler) Handle(ctx context.Context, query GetAccount) (*domain.Account, error) {
	return h.repo.Load(ctx, query.AccountID)
}
