package domain

import (
	"context"

	"github.com/stackus/edat/core"
	"github.com/stackus/edat/msg"
	"github.com/stackus/edat/saga"

	"serviceapis/orderapi"
)

func registerCreateOrderSagaData() {
	core.RegisterSagaData(CreateOrderSagaData{})
}

type CreateOrderSaga interface {
	Start(ctx context.Context, sagaData core.SagaData) (*saga.Instance, error)
	ReplyChannel() string
	ReceiveMessage(ctx context.Context, message msg.Message) error
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
