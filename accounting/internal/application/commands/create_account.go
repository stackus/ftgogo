package commands

import (
	"context"

	"github.com/stackus/edat/es"
	"github.com/stackus/ftgogo/account/internal/domain"
)

type CreateAccount struct {
	ConsumerID string
	Name       string
}

type CreateAccountHandler struct {
	repo domain.AccountRepository
}

func NewCreateAccountHandler(accountRepo domain.AccountRepository) CreateAccountHandler {
	return CreateAccountHandler{repo: accountRepo}
}

func (h CreateAccountHandler) Handle(ctx context.Context, cmd CreateAccount) error {
	_, err := h.repo.Save(ctx, &domain.CreateAccount{
		Name: cmd.Name,
	}, es.WithAggregateRootID(cmd.ConsumerID))

	return err
}
