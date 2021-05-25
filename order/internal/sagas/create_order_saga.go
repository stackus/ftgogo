package sagas

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

type CreateOrderSaga struct {
	steps []saga.Step
}

func NewCreateOrderSagaDefinition() saga.Definition {
	s := &CreateOrderSaga{}

	// TODO replicate ftgo proxies as examples
	s.steps = []saga.Step{
		saga.NewRemoteStep().
			Compensation(s.rejectOrder),
		saga.NewRemoteStep().
			Action(s.validateConsumer),
		saga.NewRemoteStep().
			Action(s.createTicket).
			HandleActionReply(kitchenapi.CreateTicketReply{}, s.handleTicketCreated).
			Compensation(s.cancelTicket),
		saga.NewRemoteStep().
			Action(s.authorizeOrder),
		saga.NewRemoteStep().
			Action(s.confirmTicket),
		saga.NewRemoteStep().
			Action(s.approveOrder),
	}

	return s
}

func (s *CreateOrderSaga) SagaName() string {
	return "orderservice.CreateOrderSaga"
}

func (s *CreateOrderSaga) ReplyChannel() string {
	return s.SagaName() + ".reply"
}

func (s *CreateOrderSaga) Steps() []saga.Step {
	return s.steps
}

func (s *CreateOrderSaga) OnHook(saga.LifecycleHook, *saga.Instance) {}

func (s *CreateOrderSaga) rejectOrder(_ context.Context, data core.SagaData) msg.DomainCommand {
	sagaData := data.(*domain.CreateOrderSagaData)

	return &orderapi.RejectOrder{
		OrderID: sagaData.OrderID,
	}
}

func (s *CreateOrderSaga) validateConsumer(_ context.Context, data core.SagaData) msg.DomainCommand {
	sagaData := data.(*domain.CreateOrderSagaData)

	return &consumerapi.ValidateOrderByConsumer{
		ConsumerID: sagaData.ConsumerID,
		OrderID:    sagaData.OrderID,
		OrderTotal: sagaData.OrderTotal,
	}
}

func (s *CreateOrderSaga) createTicket(_ context.Context, data core.SagaData) msg.DomainCommand {
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

func (s *CreateOrderSaga) handleTicketCreated(_ context.Context, data core.SagaData, reply core.Reply) error {
	sagaData := data.(*domain.CreateOrderSagaData)

	ticketCreated := reply.(*kitchenapi.CreateTicketReply)

	sagaData.TicketID = ticketCreated.TicketID

	return nil
}

func (s *CreateOrderSaga) cancelTicket(_ context.Context, data core.SagaData) msg.DomainCommand {
	sagaData := data.(*domain.CreateOrderSagaData)

	return &kitchenapi.CancelCreateTicket{
		TicketID: sagaData.TicketID,
	}
}

func (s *CreateOrderSaga) authorizeOrder(_ context.Context, data core.SagaData) msg.DomainCommand {
	sagaData := data.(*domain.CreateOrderSagaData)

	return &accountingapi.AuthorizeOrder{
		OrderID:    sagaData.OrderID,
		ConsumerID: sagaData.ConsumerID,
		OrderTotal: sagaData.OrderTotal,
	}
}

func (s *CreateOrderSaga) confirmTicket(_ context.Context, data core.SagaData) msg.DomainCommand {
	sagaData := data.(*domain.CreateOrderSagaData)

	return &kitchenapi.ConfirmCreateTicket{
		TicketID: sagaData.TicketID,
	}
}

func (s *CreateOrderSaga) approveOrder(_ context.Context, data core.SagaData) msg.DomainCommand {
	sagaData := data.(*domain.CreateOrderSagaData)

	return &orderapi.ApproveOrder{
		OrderID:  sagaData.OrderID,
		TicketID: sagaData.TicketID,
	}
}
