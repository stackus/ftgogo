package adapters

import (
	"context"

	"github.com/stackus/ftgogo/store-web/internal/domain"
)

type (
	FindAccount struct {
		AccountID string
	}

	EnableAccount struct {
		AccountID string
	}

	DisableAccount EnableAccount
)

type AccountingRepository interface {
	Find(context.Context, FindAccount) (*domain.Account, error)
	Enable(context.Context, EnableAccount) error
	Disable(context.Context, DisableAccount) error
}
