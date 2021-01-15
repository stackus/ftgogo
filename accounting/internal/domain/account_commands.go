package domain

import (
	"github.com/stackus/edat/core"
)

func registerAccountCommands() {
	core.RegisterCommands(
		CreateAccount{},
		AuthorizeOrder{}, ReverseAuthorizeOrder{}, ReviseAuthorizeOrder{},
		DisableAccount{}, EnableAccount{},
	)
}

type CreateAccount struct {
	Name string
}

func (CreateAccount) CommandName() string { return "accountingservice.CreateAccount" }

type AuthorizeOrder struct {
	OrderID    string
	OrderTotal int
}

func (AuthorizeOrder) CommandName() string { return "accountingservice.AuthorizeOrder" }

type ReverseAuthorizeOrder struct {
	OrderID    string
	OrderTotal int
}

func (ReverseAuthorizeOrder) CommandName() string { return "accountingservice.ReverseAuthorizeOrder" }

type ReviseAuthorizeOrder struct {
	OrderID    string
	OrderTotal int
}

func (ReviseAuthorizeOrder) CommandName() string { return "accountingservice.ReviseAuthorizeOrder" }

type DisableAccount struct{}

func (DisableAccount) CommandName() string { return "accountingservice.DisableAccount" }

type EnableAccount struct{}

func (EnableAccount) CommandName() string { return "accountingservice.EnableAccount" }
