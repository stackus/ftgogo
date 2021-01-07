package accountingapi

import (
	"github.com/stackus/edat/core"
)

func registerCommands() {
	core.RegisterCommands(AuthorizeOrder{}, ReverseAuthorizeOrder{}, ReviseAuthorizeOrder{})
}

type AccountServiceCommand struct{}

func (AccountServiceCommand) DestinationChannel() string { return AccountingServiceCommandChannel }

type AuthorizeOrder struct {
	AccountServiceCommand
	ConsumerID string
	OrderID    string
	OrderTotal int // Money
}

func (AuthorizeOrder) CommandName() string { return "accountingapi.AuthorizeOrder" }

type ReverseAuthorizeOrder struct {
	AccountServiceCommand
	ConsumerID string
	OrderID    string
	OrderTotal int // Money
}

func (ReverseAuthorizeOrder) CommandName() string { return "accountingapi.ReverseAuthorizeOrder" }

type ReviseAuthorizeOrder struct {
	AccountServiceCommand
	ConsumerID string
	OrderID    string
	OrderTotal int // Money
}

func (ReviseAuthorizeOrder) CommandName() string { return "accountingapi.ReviseAuthorizeOrder" }
