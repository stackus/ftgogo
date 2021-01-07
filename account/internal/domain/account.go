package domain

import (
	"fmt"

	"github.com/stackus/edat/core"
	"github.com/stackus/edat/es"
	"shared-go/errs"
)

var ErrAccountUnhandledCommand = errs.NewError("unhandled command in account aggregate", errs.ErrServerError)
var ErrAccountUnhandledEvent = errs.NewError("unhandled event in account aggregate", errs.ErrServerError)
var ErrAccountUnhandledSnapshot = errs.NewError("unhandled snapshot in account aggregate", errs.ErrServerError)

var ErrAccountDisabled = errs.NewError("account disabled", errs.ErrBadRequest)
var ErrAccountEnabled = errs.NewError("account enabled", errs.ErrBadRequest)

type Account struct {
	es.AggregateBase
	Name    string
	Enabled bool
}

var _ es.Aggregate = (*Account)(nil)

func NewAccount() es.Aggregate {
	return &Account{}
}

func (a *Account) EntityName() string {
	return "accountingservice.Account"
}

func (a *Account) ProcessCommand(command core.Command) error {
	switch cmd := command.(type) {
	case *CreateAccount:
		a.AddEvent(
			&AccountCreated{Name: cmd.Name},
			&AccountEnabled{},
		)

	case *AuthorizeOrder:
		if !a.Enabled {
			return ErrAccountDisabled
		}

		a.AddEvent(&OrderAuthorized{
			OrderID:    cmd.OrderID,
			OrderTotal: cmd.OrderTotal,
		})

	case *ReverseAuthorizeOrder:
		if !a.Enabled {
			return ErrAccountDisabled
		}

		// noop

	case *ReviseAuthorizeOrder:
		if !a.Enabled {
			return ErrAccountDisabled
		}

		// noop

	case *DisableAccount:
		if !a.Enabled {
			return ErrAccountDisabled
		}

		a.AddEvent(&AccountDisabled{})

	case *EnableAccount:
		if a.Enabled {
			return ErrAccountEnabled
		}

		a.AddEvent(&AccountEnabled{})

	default:
		return errs.NewError(fmt.Sprintf("unhandled command `%T`", command), ErrAccountUnhandledCommand)
	}

	return nil
}

// ApplyEvent makes changes to the aggregate based on the event and its payload
func (a *Account) ApplyEvent(event core.Event) error {
	switch evt := event.(type) {
	case *AccountCreated:
		a.Name = evt.Name

	case *OrderAuthorized:

	case *AccountEnabled:
		a.Enabled = true

	case *AccountDisabled:
		a.Enabled = false

	default:
		return errs.NewError(fmt.Sprintf("unhandled event `%T`", event), ErrAccountUnhandledEvent)
	}

	return nil
}

func (a *Account) ApplySnapshot(snapshot core.Snapshot) error {
	switch ss := snapshot.(type) {
	case *AccountSnapshot:
		a.Name = ss.Name
		a.Enabled = ss.Enabled

	default:
		return errs.NewError(fmt.Sprintf("unhandled snapshot `%T`", snapshot), ErrAccountUnhandledSnapshot)
	}

	return nil
}

func (a *Account) ToSnapshot() (core.Snapshot, error) {
	return &AccountSnapshot{
		Name:    a.Name,
		Enabled: a.Enabled,
	}, nil
}
