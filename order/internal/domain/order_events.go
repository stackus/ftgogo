package domain

import (
	"github.com/stackus/edat/core"
)

func registerOrderEvents() {
	core.RegisterEvents(
		OrderCancelling{}, OrderCancellingUndone{},
		OrderRevisingUndone{},
	)
}

type OrderCancelling struct{}

// EventName event method
func (OrderCancelling) EventName() string { return "orderservice.OrderCancelling" }

type OrderCancellingUndone struct{}

// EventName event method
func (OrderCancellingUndone) EventName() string { return "orderservice.OrderCancellingUndone" }

type OrderRevisingUndone struct{}

// EventName event method
func (OrderRevisingUndone) EventName() string { return "orderservice.OrderRevisingUndone" }
