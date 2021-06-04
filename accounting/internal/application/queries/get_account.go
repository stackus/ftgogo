package queries

import (
	"context"

	"github.com/stackus/ftgogo/accounting/internal/adapters"
	"github.com/stackus/ftgogo/accounting/internal/domain"
)

type GetAccount struct {
	AccountID string
}

type GetAccountHandler struct {
	repo adapters.AccountRepository
}

func NewGetAccountHandler(accountRepo adapters.AccountRepository) GetAccountHandler {
	return GetAccountHandler{repo: accountRepo}
}

func (h GetAccountHandler) Handle(ctx context.Context, query GetAccount) (*domain.Account, error) {
	root, err := h.repo.Load(ctx, query.AccountID)
	if err != nil {
		return nil, err
	}

	return root.Aggregate().(*domain.Account), nil
}
