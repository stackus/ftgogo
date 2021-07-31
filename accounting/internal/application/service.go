package application

import (
	"context"

	"github.com/stackus/ftgogo/accounting/internal/application/commands"
	"github.com/stackus/ftgogo/accounting/internal/application/ports"
	"github.com/stackus/ftgogo/accounting/internal/application/queries"
	"github.com/stackus/ftgogo/accounting/internal/domain"
)

type ServiceApplication interface {
	AuthorizeOrder(ctx context.Context, cmd commands.AuthorizeOrder) error
	CreateAccount(ctx context.Context, cmd commands.CreateAccount) error
	DisableAccount(ctx context.Context, cmd commands.DisableAccount) error
	EnableAccount(ctx context.Context, cmd commands.EnableAccount) error
	GetAccount(ctx context.Context, query queries.GetAccount) (*domain.Account, error)
	ReverseAuthorizeOrder(ctx context.Context, cmd commands.ReverseAuthorizeOrder) error
	ReviseAuthorizeOrder(ctx context.Context, cmd commands.ReviseAuthorizeOrder) error
}

type Service struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	AuthorizeOrder        commands.AuthorizeOrderHandler
	CreateAccount         commands.CreateAccountHandler
	DisableAccount        commands.DisableAccountHandler
	EnableAccount         commands.EnableAccountHandler
	ReverseAuthorizeOrder commands.ReverseAuthorizeOrderHandler
	ReviseAuthorizeOrder  commands.ReviseAuthorizeOrderHandler
}

type Queries struct {
	GetAccount queries.GetAccountHandler
}

func NewServiceApplication(accountRepo ports.AccountRepository) *Service {
	return &Service{
		Commands: Commands{
			AuthorizeOrder:        commands.NewAuthorizeOrderHandler(accountRepo),
			CreateAccount:         commands.NewCreateAccountHandler(accountRepo),
			DisableAccount:        commands.NewDisableAccountHandler(accountRepo),
			EnableAccount:         commands.NewEnableAccountHandler(accountRepo),
			ReverseAuthorizeOrder: commands.NewReverseAuthorizeOrderHandler(accountRepo),
			ReviseAuthorizeOrder:  commands.NewReviseAuthorizeOrderHandler(accountRepo),
		},
		Queries: Queries{
			GetAccount: queries.NewGetAccountHandler(accountRepo),
		},
	}
}

func (s Service) AuthorizeOrder(ctx context.Context, cmd commands.AuthorizeOrder) error {
	return s.Commands.AuthorizeOrder.Handle(ctx, cmd)
}

func (s Service) CreateAccount(ctx context.Context, cmd commands.CreateAccount) error {
	return s.Commands.CreateAccount.Handle(ctx, cmd)
}

func (s Service) DisableAccount(ctx context.Context, cmd commands.DisableAccount) error {
	return s.Commands.DisableAccount.Handle(ctx, cmd)
}

func (s Service) EnableAccount(ctx context.Context, cmd commands.EnableAccount) error {
	return s.Commands.EnableAccount.Handle(ctx, cmd)
}

func (s Service) GetAccount(ctx context.Context, query queries.GetAccount) (*domain.Account, error) {
	return s.Queries.GetAccount.Handle(ctx, query)
}

func (s Service) ReverseAuthorizeOrder(ctx context.Context, cmd commands.ReverseAuthorizeOrder) error {
	return s.Commands.ReverseAuthorizeOrder.Handle(ctx, cmd)
}

func (s Service) ReviseAuthorizeOrder(ctx context.Context, cmd commands.ReviseAuthorizeOrder) error {
	return s.Commands.ReviseAuthorizeOrder.Handle(ctx, cmd)
}
