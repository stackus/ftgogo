package domain

import (
	"context"

	"github.com/stackus/edat/core"
	"github.com/stackus/edat/msg"
	"github.com/stackus/edat/saga"
	"serviceapis/commonapi"
)

func registerReviseOrderSagaData() {
	core.RegisterSagaData(ReviseOrderSagaData{})
}

type ReviseOrderSaga interface {
	Start(ctx context.Context, sagaData core.SagaData) (*saga.Instance, error)
	ReplyChannel() string
	ReceiveMessage(ctx context.Context, message msg.Message) error
}

type ReviseOrderSagaData struct {
	OrderID           string
	ConsumerID        string
	RestaurantID      string
	TicketID          string
	ExpectedVersion   int
	RevisedOrderTotal int
	RevisedQuantities commonapi.MenuItemQuantities
}

func (ReviseOrderSagaData) SagaDataName() string {
	return "orderservice.ReviseOrderSagaData"
}
