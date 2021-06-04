package domain

import (
	"github.com/stackus/edat/core"
)

func registerCancelOrderSagaData() {
	core.RegisterSagaData(CancelOrderSagaData{})
}

type CancelOrderSagaData struct {
	OrderID      string
	ConsumerID   string
	RestaurantID string
	TicketID     string
	OrderTotal   int
}

func (CancelOrderSagaData) SagaDataName() string {
	return "orderservice.CancelOrderSagaData"
}
