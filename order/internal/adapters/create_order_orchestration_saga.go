package adapters

import (
	"context"

	"github.com/stackus/edat/core"
	"github.com/stackus/edat/msg"
	"github.com/stackus/edat/saga"

	"github.com/stackus/ftgogo/order/internal/domain"
	"github.com/stackus/ftgogo/serviceapis/accountingapi"
	"github.com/stackus/ftgogo/serviceapis/consumerapi"
	"github.com/stackus/ftgogo/serviceapis/kitchenapi"
	"github.com/stackus/ftgogo/serviceapis/orderapi"
)

type createOrderSagaDefinition struct {
	steps []saga.Step
}

func NewCreateOrderOrchestrationSaga(store saga.InstanceStore, publisher msg.CommandMessagePublisher, options ...saga.OrchestratorOption) *saga.Orchestrator {
	definition := &createOrderSagaDefinition{}

	// TODO replicate ftgo proxies as examples
	definition.steps = []saga.Step{
		saga.NewRemoteStep().
			Compensation(definition.rejectOrder),
		saga.NewRemoteStep().
			Action(definition.validateConsumer),
		saga.NewRemoteStep().
			Action(definition.createTicket).
			HandleActionReply(kitchenapi.CreateTicketReply{}, definition.handleTicketCreated).
			Compensation(definition.cancelTicket),
		saga.NewRemoteStep().
			Action(definition.authorizeOrder),
		saga.NewRemoteStep().
			Action(definition.confirmTicket),
		saga.NewRemoteStep().
			Action(definition.approveOrder),
	}

	return saga.NewOrchestrator(definition, store, publisher, options...)
}

func (s *createOrderSagaDefinition) SagaName() string {
	return "orderservice.CreateOrderSaga"
}

func (s *createOrderSagaDefinition) ReplyChannel() string {
	return s.SagaName() + ".reply"
}

func (s *createOrderSagaDefinition) Steps() []saga.Step {
	return s.steps
}

func (s *createOrderSagaDefinition) OnHook(saga.LifecycleHook, *saga.Instance) {}

func (s *createOrderSagaDefinition) rejectOrder(_ context.Context, data core.SagaData) msg.DomainCommand {
	sagaData := data.(*domain.CreateOrderSagaData)

	return &orderapi.RejectOrder{
		OrderID: sagaData.OrderID,
	}
}

func (s *createOrderSagaDefinition) validateConsumer(_ context.Context, data core.SagaData) msg.DomainCommand {
	sagaData := data.(*domain.CreateOrderSagaData)

	return &consumerapi.ValidateOrderByConsumer{
		ConsumerID: sagaData.ConsumerID,
		OrderID:    sagaData.OrderID,
		OrderTotal: sagaData.OrderTotal,
	}
}

func (s *createOrderSagaDefinition) createTicket(_ context.Context, data core.SagaData) msg.DomainCommand {
	sagaData := data.(*domain.CreateOrderSagaData)

	ticketDetails := make([]kitchenapi.LineItem, 0)
	for _, lineItem := range sagaData.LineItems {
		ticketDetails = append(ticketDetails, kitchenapi.LineItem{
			Name:       lineItem.Name,
			MenuItemID: lineItem.MenuItemID,
			Quantity:   lineItem.Quantity,
		})
	}

	return &kitchenapi.CreateTicket{
		OrderID:       sagaData.OrderID,
		RestaurantID:  sagaData.RestaurantID,
		TicketDetails: ticketDetails,
	}
}

func (s *createOrderSagaDefinition) handleTicketCreated(_ context.Context, data core.SagaData, reply core.Reply) error {
	sagaData := data.(*domain.CreateOrderSagaData)

	ticketCreated := reply.(*kitchenapi.CreateTicketReply)

	sagaData.TicketID = ticketCreated.TicketID

	return nil
}

func (s *createOrderSagaDefinition) cancelTicket(_ context.Context, data core.SagaData) msg.DomainCommand {
	sagaData := data.(*domain.CreateOrderSagaData)

	return &kitchenapi.CancelCreateTicket{
		TicketID: sagaData.TicketID,
	}
}

func (s *createOrderSagaDefinition) authorizeOrder(_ context.Context, data core.SagaData) msg.DomainCommand {
	sagaData := data.(*domain.CreateOrderSagaData)

	return &accountingapi.AuthorizeOrder{
		OrderID:    sagaData.OrderID,
		ConsumerID: sagaData.ConsumerID,
		OrderTotal: sagaData.OrderTotal,
	}
}

func (s *createOrderSagaDefinition) confirmTicket(_ context.Context, data core.SagaData) msg.DomainCommand {
	sagaData := data.(*domain.CreateOrderSagaData)

	return &kitchenapi.ConfirmCreateTicket{
		TicketID: sagaData.TicketID,
	}
}

func (s *createOrderSagaDefinition) approveOrder(_ context.Context, data core.SagaData) msg.DomainCommand {
	sagaData := data.(*domain.CreateOrderSagaData)

	return &orderapi.ApproveOrder{
		OrderID:  sagaData.OrderID,
		TicketID: sagaData.TicketID,
	}
}
