package application

import (
	"github.com/stackus/ftgogo/order/internal/application/commands"
	"github.com/stackus/ftgogo/order/internal/application/queries"
)

type Service struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateOrder          commands.CreateOrderHandler
	ApproveOrder         commands.ApproveOrderHandler
	RejectOrder          commands.RejectOrderHandler
	BeginCancelOrder     commands.BeginCancelOrderHandler
	UndoCancelOrder      commands.UndoCancelOrderHandler
	ConfirmCancelOrder   commands.ConfirmCancelOrderHandler
	BeginReviseOrder     commands.BeginReviseOrderHandler
	UndoReviseOrder      commands.UndoReviseOrderHandler
	ConfirmReviseOrder   commands.ConfirmReviseOrderHandler
	StartCreateOrderSaga commands.StartCreateOrderSagaHandler
	StartCancelOrderSaga commands.StartCancelOrderSagaHandler
	StartReviseOrderSaga commands.StartReviseOrderSagaHandler
	CreateRestaurant     commands.CreateRestaurantHandler
	ReviseRestaurantMenu commands.ReviseRestaurantMenuHandler
}

type Queries struct {
	GetOrder      queries.GetOrderHandler
	GetRestaurant queries.GetRestaurantHandler
}
