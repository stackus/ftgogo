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

type reviseOrderSagaDefinition struct {
	steps []saga.Step
}

func NewReviseOrderOrchestrationSaga(store saga.InstanceStore, publisher msg.CommandMessagePublisher, options ...saga.OrchestratorOption) *saga.Orchestrator {
	definition := &reviseOrderSagaDefinition{}

	definition.steps = []saga.Step{
		saga.NewRemoteStep().
			Action(definition.beginReviseOrder).
			HandleActionReply(orderapi.BeginReviseOrderReply{}, definition.handleBeginReviseOrderReply).
			Compensation(definition.undoReviseOrder),
		saga.NewRemoteStep().
			Action(definition.beginReviseTicket).
			Compensation(definition.undoReviseTicket),
		saga.NewRemoteStep().
			Action(definition.reviseAuthorization),
		saga.NewRemoteStep().
			Action(definition.confirmTicketRevision),
		saga.NewRemoteStep().
			Action(definition.confirmOrderRevision),
	}

	return saga.NewOrchestrator(definition, store, publisher, options...)
}

func (s *reviseOrderSagaDefinition) SagaName() string {
	return "orderservice.ReviseOrderSaga"
}

func (s *reviseOrderSagaDefinition) ReplyChannel() string {
	return s.SagaName() + ".reply"
}

func (s *reviseOrderSagaDefinition) Steps() []saga.Step {
	return s.steps
}

func (s *reviseOrderSagaDefinition) OnHook(saga.LifecycleHook, *saga.Instance) {}

func (s *reviseOrderSagaDefinition) beginReviseOrder(_ context.Context, data core.SagaData) msg.DomainCommand {
	sagaData := data.(*domain.ReviseOrderSagaData)

	return &orderapi.BeginReviseOrder{
		OrderID:           sagaData.OrderID,
		RevisedQuantities: sagaData.RevisedQuantities,
	}
}

func (s *reviseOrderSagaDefinition) undoReviseOrder(_ context.Context, data core.SagaData) msg.DomainCommand {
	sagaData := data.(*domain.ReviseOrderSagaData)

	return &orderapi.UndoReviseOrder{OrderID: sagaData.OrderID}
}

func (s *reviseOrderSagaDefinition) handleBeginReviseOrderReply(_ context.Context, data core.SagaData, reply core.Reply) error {
	sagaData := data.(*domain.ReviseOrderSagaData)
	reviseOrderReply := reply.(*orderapi.BeginReviseOrderReply)

	sagaData.RevisedOrderTotal = reviseOrderReply.RevisedOrderTotal

	return nil
}

func (s *reviseOrderSagaDefinition) beginReviseTicket(_ context.Context, data core.SagaData) msg.DomainCommand {
	sagaData := data.(*domain.ReviseOrderSagaData)

	return &kitchenapi.BeginReviseTicket{
		TicketID:          sagaData.TicketID,
		RestaurantID:      sagaData.RestaurantID,
		RevisedQuantities: sagaData.RevisedQuantities,
	}
}

func (s *reviseOrderSagaDefinition) undoReviseTicket(_ context.Context, data core.SagaData) msg.DomainCommand {
	sagaData := data.(*domain.ReviseOrderSagaData)

	return &kitchenapi.UndoReviseTicket{
		TicketID:     sagaData.TicketID,
		RestaurantID: sagaData.RestaurantID,
	}
}

func (s *reviseOrderSagaDefinition) reviseAuthorization(_ context.Context, data core.SagaData) msg.DomainCommand {
	sagaData := data.(*domain.ReviseOrderSagaData)

	return &accountingapi.ReviseAuthorizeOrder{
		ConsumerID: sagaData.ConsumerID,
		OrderID:    sagaData.OrderID,
		OrderTotal: sagaData.RevisedOrderTotal,
	}
}

func (s *reviseOrderSagaDefinition) confirmTicketRevision(_ context.Context, data core.SagaData) msg.DomainCommand {
	sagaData := data.(*domain.ReviseOrderSagaData)

	return &kitchenapi.ConfirmReviseTicket{
		TicketID:          sagaData.TicketID,
		RestaurantID:      sagaData.RestaurantID,
		RevisedQuantities: sagaData.RevisedQuantities,
	}
}

func (s *reviseOrderSagaDefinition) confirmOrderRevision(_ context.Context, data core.SagaData) msg.DomainCommand {
	sagaData := data.(*domain.ReviseOrderSagaData)

	return &orderapi.ConfirmReviseOrder{OrderID: sagaData.OrderID}
}
