package commands

import (
	"context"

	"github.com/stackus/edat/es"
	"github.com/stackus/errors"

	"github.com/stackus/ftgogo/accounting/internal/application/ports"
	"github.com/stackus/ftgogo/accounting/internal/domain"
)

type CreateAccount struct {
	ConsumerID string
	Name       string
}

type CreateAccountHandler struct {
	repo ports.AccountRepository
}

func NewCreateAccountHandler(accountRepo ports.AccountRepository) CreateAccountHandler {
	return CreateAccountHandler{repo: accountRepo}
}

func (h CreateAccountHandler) Handle(ctx context.Context, cmd CreateAccount) error {
	_, err := h.repo.Save(ctx, &domain.CreateAccount{
		Name: cmd.Name,
	}, es.WithAggregateRootID(cmd.ConsumerID))

	// TODO update edat-pgx to return an es.ErrVersionConflict when versions do not align
	// for now assume all errors are duplicate account errors
	if err != nil {
		return errors.Wrap(errors.ErrConflict, "account already exists")
	}

	return nil
}
