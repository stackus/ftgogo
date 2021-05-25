package sagas

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

type ReviseOrderSaga struct {
	steps []saga.Step
}

func NewReviseOrderSagaDefinition() saga.Definition {
	s := &ReviseOrderSaga{}

	s.steps = []saga.Step{
		saga.NewRemoteStep().
			Action(s.beginReviseOrder).
			HandleActionReply(orderapi.BeginReviseOrderReply{}, s.handleBeginReviseOrderReply).
			Compensation(s.undoReviseOrder),
		saga.NewRemoteStep().
			Action(s.beginReviseTicket).
			Compensation(s.undoReviseTicket),
		saga.NewRemoteStep().
			Action(s.reviseAuthorization),
		saga.NewRemoteStep().
			Action(s.confirmTicketRevision),
		saga.NewRemoteStep().
			Action(s.confirmOrderRevision),
	}

	return s
}

func (s *ReviseOrderSaga) SagaName() string {
	return "orderservice.ReviseOrderSaga"
}

func (s *ReviseOrderSaga) ReplyChannel() string {
	return s.SagaName() + ".reply"
}

func (s *ReviseOrderSaga) Steps() []saga.Step {
	return s.steps
}

func (s *ReviseOrderSaga) OnHook(saga.LifecycleHook, *saga.Instance) {}

func (s *ReviseOrderSaga) beginReviseOrder(_ context.Context, data core.SagaData) msg.DomainCommand {
	sagaData := data.(*domain.ReviseOrderSagaData)

	return &orderapi.BeginReviseOrder{
		OrderID:           sagaData.OrderID,
		RevisedQuantities: sagaData.RevisedQuantities,
	}
}

func (s *ReviseOrderSaga) undoReviseOrder(_ context.Context, data core.SagaData) msg.DomainCommand {
	sagaData := data.(*domain.ReviseOrderSagaData)

	return &orderapi.UndoReviseOrder{OrderID: sagaData.OrderID}
}

func (s *ReviseOrderSaga) handleBeginReviseOrderReply(_ context.Context, data core.SagaData, reply core.Reply) error {
	sagaData := data.(*domain.ReviseOrderSagaData)
	reviseOrderReply := reply.(*orderapi.BeginReviseOrderReply)

	sagaData.RevisedOrderTotal = reviseOrderReply.RevisedOrderTotal

	return nil
}

func (s *ReviseOrderSaga) beginReviseTicket(_ context.Context, data core.SagaData) msg.DomainCommand {
	sagaData := data.(*domain.ReviseOrderSagaData)

	return &kitchenapi.BeginReviseTicket{
		TicketID:          sagaData.TicketID,
		RestaurantID:      sagaData.RestaurantID,
		RevisedQuantities: sagaData.RevisedQuantities,
	}
}

func (s *ReviseOrderSaga) undoReviseTicket(_ context.Context, data core.SagaData) msg.DomainCommand {
	sagaData := data.(*domain.ReviseOrderSagaData)

	return &kitchenapi.UndoReviseTicket{
		TicketID:     sagaData.TicketID,
		RestaurantID: sagaData.RestaurantID,
	}
}

func (s *ReviseOrderSaga) reviseAuthorization(_ context.Context, data core.SagaData) msg.DomainCommand {
	sagaData := data.(*domain.ReviseOrderSagaData)

	return &accountingapi.ReviseAuthorizeOrder{
		ConsumerID: sagaData.ConsumerID,
		OrderID:    sagaData.OrderID,
		OrderTotal: sagaData.RevisedOrderTotal,
	}
}

func (s *ReviseOrderSaga) confirmTicketRevision(_ context.Context, data core.SagaData) msg.DomainCommand {
	sagaData := data.(*domain.ReviseOrderSagaData)

	return &kitchenapi.ConfirmReviseTicket{
		TicketID:          sagaData.TicketID,
		RestaurantID:      sagaData.RestaurantID,
		RevisedQuantities: sagaData.RevisedQuantities,
	}
}

func (s *ReviseOrderSaga) confirmOrderRevision(_ context.Context, data core.SagaData) msg.DomainCommand {
	sagaData := data.(*domain.ReviseOrderSagaData)

	return &orderapi.ConfirmReviseOrder{OrderID: sagaData.OrderID}
}
