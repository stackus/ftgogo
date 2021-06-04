package application

import (
	"github.com/stackus/ftgogo/accounting/internal/application/commands"
	"github.com/stackus/ftgogo/accounting/internal/application/queries"
)

type Service struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	AuthorizeOrder        commands.AuthorizeOrderHandler
	ReverseAuthorizeOrder commands.ReverseAuthorizeOrderHandler
	ReviseAuthorizeOrder  commands.ReviseAuthorizeOrderHandler
	CreateAccount         commands.CreateAccountHandler
	DisableAccount        commands.DisableAccountHandler
	EnableAccount         commands.EnableAccountHandler
}

type Queries struct {
	GetAccount queries.GetAccountHandler
}
