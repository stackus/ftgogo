package domain

import (
	"context"

	"github.com/stackus/edat/core"
	"github.com/stackus/edat/msg"
	"github.com/stackus/edat/saga"
)

func registerCancelOrderSagaData() {
	core.RegisterSagaData(CancelOrderSagaData{})
}

type CancelOrderSaga interface {
	Start(ctx context.Context, sagaData core.SagaData) (*saga.Instance, error)
	ReplyChannel() string
	ReceiveMessage(ctx context.Context, message msg.Message) error
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
