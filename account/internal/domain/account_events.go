package domain

import (
	"github.com/stackus/edat/core"
)

func registerAccountEvents() {
	core.RegisterEvents(
		AccountCreated{}, OrderAuthorized{},
		AccountDisabled{}, AccountEnabled{},
	)
}

type AccountCreated struct {
	Name string
}

func (AccountCreated) EventName() string { return "accountingservice.AccountCreated" }

type OrderAuthorized struct {
	OrderID    string
	OrderTotal int
}

func (OrderAuthorized) EventName() string { return "accountingservice.OrderAuthorized" }

type AccountDisabled struct{}

func (AccountDisabled) EventName() string { return "accountingservice.AccountDisabled" }

type AccountEnabled struct{}

func (AccountEnabled) EventName() string { return "accountingservice.AccountEnabled" }
