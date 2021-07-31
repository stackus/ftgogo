package adapters

import (
	"context"

	"github.com/stackus/edat/core"
	"github.com/stackus/edat/msg"
	"github.com/stackus/edat/saga"

	"github.com/stackus/ftgogo/order/internal/domain"
	"github.com/stackus/ftgogo/serviceapis/accountingapi"
	"github.com/stackus/ftgogo/serviceapis/kitchenapi"
	"github.com/stackus/ftgogo/serviceapis/orderapi"
)

type cancelOrderSagaDefinition struct {
	steps []saga.Step
}

func NewCancelOrderOrchestrationSaga(store saga.InstanceStore, publisher msg.CommandMessagePublisher, options ...saga.OrchestratorOption) *saga.Orchestrator {
	definition := &cancelOrderSagaDefinition{}

	definition.steps = []saga.Step{
		saga.NewRemoteStep().
			Action(definition.beginCancelOrder).
			Compensation(definition.undoBeginCancelOrder),
		saga.NewRemoteStep().
			Action(definition.beginCancelTicket).
			Compensation(definition.undoBeginCancelTicket),
		saga.NewRemoteStep().
			Action(definition.reverseAuthorization),
		saga.NewRemoteStep().
			Action(definition.confirmCancelTicket),
		saga.NewRemoteStep().
			Action(definition.confirmCancelOrder),
	}

	return saga.NewOrchestrator(definition, store, publisher, options...)
}

func (s *cancelOrderSagaDefinition) SagaName() string {
	return "orderservice.CancelOrderSaga"
}

func (s *cancelOrderSagaDefinition) ReplyChannel() string {
	return s.SagaName() + ".reply"
}

func (s *cancelOrderSagaDefinition) Steps() []saga.Step {
	return s.steps
}

func (s *cancelOrderSagaDefinition) OnHook(saga.LifecycleHook, *saga.Instance) {}

func (s *cancelOrderSagaDefinition) beginCancelOrder(_ context.Context, data core.SagaData) msg.DomainCommand {
	sagaData := data.(*domain.CancelOrderSagaData)

	return &orderapi.BeginCancelOrder{OrderID: sagaData.OrderID}
}

func (s *cancelOrderSagaDefinition) undoBeginCancelOrder(_ context.Context, data core.SagaData) msg.DomainCommand {
	sagaData := data.(*domain.CancelOrderSagaData)

	return &orderapi.UndoCancelOrder{OrderID: sagaData.OrderID}
}

func (s *cancelOrderSagaDefinition) beginCancelTicket(_ context.Context, data core.SagaData) msg.DomainCommand {
	sagaData := data.(*domain.CancelOrderSagaData)

	return &kitchenapi.BeginCancelTicket{
		TicketID:     sagaData.TicketID,
		RestaurantID: sagaData.RestaurantID,
	}
}

func (s *cancelOrderSagaDefinition) undoBeginCancelTicket(_ context.Context, data core.SagaData) msg.DomainCommand {
	sagaData := data.(*domain.CancelOrderSagaData)

	return &kitchenapi.UndoCancelTicket{
		TicketID:     sagaData.TicketID,
		RestaurantID: sagaData.RestaurantID,
	}
}

func (s *cancelOrderSagaDefinition) reverseAuthorization(_ context.Context, data core.SagaData) msg.DomainCommand {
	sagaData := data.(*domain.CancelOrderSagaData)

	return &accountingapi.ReverseAuthorizeOrder{
		ConsumerID: sagaData.ConsumerID,
		OrderID:    sagaData.OrderID,
		OrderTotal: sagaData.OrderTotal,
	}
}

func (s *cancelOrderSagaDefinition) confirmCancelTicket(_ context.Context, data core.SagaData) msg.DomainCommand {
	sagaData := data.(*domain.CancelOrderSagaData)

	return &kitchenapi.ConfirmCancelTicket{TicketID: sagaData.TicketID}
}

func (s *cancelOrderSagaDefinition) confirmCancelOrder(_ context.Context, data core.SagaData) msg.DomainCommand {
	sagaData := data.(*domain.CancelOrderSagaData)

	return &orderapi.ConfirmCancelOrder{OrderID: sagaData.OrderID}
}
