package domain

import (
	"github.com/stackus/edat/core"

	"github.com/stackus/ftgogo/serviceapis/orderapi"
)

func registerCreateOrderSagaData() {
	core.RegisterSagaData(CreateOrderSagaData{})
}

type CreateOrderSagaData struct {
	OrderID      string
	TicketID     string
	ConsumerID   string
	RestaurantID string
	LineItems    []orderapi.LineItem
	OrderTotal   int
}

func (CreateOrderSagaData) SagaDataName() string {
	return "orderservice.CreateOrderSagaData"
}
