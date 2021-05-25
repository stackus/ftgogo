package sagas

import (
	"context"

	"github.com/stackus/edat/core"
	"github.com/stackus/edat/msg"
	"github.com/stackus/edat/saga"

	"github.com/stackus/ftgogo/serviceapis/accountingapi"
	"github.com/stackus/ftgogo/serviceapis/kitchenapi"
	"github.com/stackus/ftgogo/serviceapis/orderapi"

	"github.com/stackus/ftgogo/order/internal/domain"
)

type CancelOrderSaga struct {
	steps []saga.Step
}

func NewCancelOrderSagaDefinition() saga.Definition {
	s := &CancelOrderSaga{}

	s.steps = []saga.Step{
		saga.NewRemoteStep().
			Action(s.beginCancelOrder).
			Compensation(s.undoBeginCancelOrder),
		saga.NewRemoteStep().
			Action(s.beginCancelTicket).
			Compensation(s.undoBeginCancelTicket),
		saga.NewRemoteStep().
			Action(s.reverseAuthorization),
		saga.NewRemoteStep().
			Action(s.confirmCancelTicket),
		saga.NewRemoteStep().
			Action(s.confirmCancelOrder),
	}

	return s
}

func (s *CancelOrderSaga) SagaName() string {
	return "orderservice.CancelOrderSaga"
}

func (s *CancelOrderSaga) ReplyChannel() string {
	return s.SagaName() + ".reply"
}

func (s *CancelOrderSaga) Steps() []saga.Step {
	return s.steps
}

func (s *CancelOrderSaga) OnHook(saga.LifecycleHook, *saga.Instance) {}

func (s *CancelOrderSaga) beginCancelOrder(_ context.Context, data core.SagaData) msg.DomainCommand {
	sagaData := data.(*domain.CancelOrderSagaData)

	return &orderapi.BeginCancelOrder{OrderID: sagaData.OrderID}
}

func (s *CancelOrderSaga) undoBeginCancelOrder(_ context.Context, data core.SagaData) msg.DomainCommand {
	sagaData := data.(*domain.CancelOrderSagaData)

	return &orderapi.UndoCancelOrder{OrderID: sagaData.OrderID}
}

func (s *CancelOrderSaga) beginCancelTicket(_ context.Context, data core.SagaData) msg.DomainCommand {
	sagaData := data.(*domain.CancelOrderSagaData)

	return &kitchenapi.BeginCancelTicket{
		TicketID:     sagaData.TicketID,
		RestaurantID: sagaData.RestaurantID,
	}
}

func (s *CancelOrderSaga) undoBeginCancelTicket(_ context.Context, data core.SagaData) msg.DomainCommand {
	sagaData := data.(*domain.CancelOrderSagaData)

	return &kitchenapi.UndoCancelTicket{
		TicketID:     sagaData.TicketID,
		RestaurantID: sagaData.RestaurantID,
	}
}

func (s *CancelOrderSaga) reverseAuthorization(_ context.Context, data core.SagaData) msg.DomainCommand {
	sagaData := data.(*domain.CancelOrderSagaData)

	return &accountingapi.ReverseAuthorizeOrder{
		ConsumerID: sagaData.ConsumerID,
		OrderID:    sagaData.OrderID,
		OrderTotal: sagaData.OrderTotal,
	}
}

func (s *CancelOrderSaga) confirmCancelTicket(_ context.Context, data core.SagaData) msg.DomainCommand {
	sagaData := data.(*domain.CancelOrderSagaData)

	return &kitchenapi.ConfirmCancelTicket{TicketID: sagaData.TicketID}
}

func (s *CancelOrderSaga) confirmCancelOrder(_ context.Context, data core.SagaData) msg.DomainCommand {
	sagaData := data.(*domain.CancelOrderSagaData)

	return &orderapi.ConfirmCancelOrder{OrderID: sagaData.OrderID}
}
