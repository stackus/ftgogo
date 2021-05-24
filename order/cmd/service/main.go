package main

import (
	"github.com/stackus/edat/msg"
	"github.com/stackus/edat/saga"

	"github.com/stackus/ftgogo/order/internal/adapters"
	"github.com/stackus/ftgogo/order/internal/application/commands"
	"github.com/stackus/ftgogo/order/internal/application/queries"
	"github.com/stackus/ftgogo/order/internal/domain"
	"github.com/stackus/ftgogo/serviceapis"
	"github.com/stackus/ftgogo/serviceapis/orderapi"
	"github.com/stackus/ftgogo/serviceapis/orderapi/pb"
	"github.com/stackus/ftgogo/serviceapis/restaurantapi"
	"shared-go/applications"
)

type Application struct {
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
	GetOrder queries.GetOrderHandler
}

func main() {
	svc := applications.NewService(initService)
	if err := svc.Execute(); err != nil {
		panic(err)
	}
}

func initService(svc *applications.Service) error {
	serviceapis.RegisterTypes()
	domain.RegisterTypes()

	orderRepo := adapters.NewOrderRepositoryPublisherMiddleware(
		adapters.NewOrderRepository(svc.AggregateStore),
		adapters.NewOrderPublisher(svc.Publisher),
	)
	restaurantRepo := adapters.NewRestaurantPostgresRepository(svc.PgConn)

	createOrderSaga := adapters.NewCreateOrderSaga(svc.SagaInstanceStore, svc.Publisher)
	svc.Subscriber.Subscribe(createOrderSaga.ReplyChannel(), createOrderSaga)

	cancelOrderSaga := adapters.NewCancelOrderSaga(svc.SagaInstanceStore, svc.Publisher)
	svc.Subscriber.Subscribe(cancelOrderSaga.ReplyChannel(), cancelOrderSaga)

	reviseOrderSaga := adapters.NewReviseOrderSaga(svc.SagaInstanceStore, svc.Publisher)
	svc.Subscriber.Subscribe(reviseOrderSaga.ReplyChannel(), reviseOrderSaga)

	ordersPlacedCounter := adapters.NewOrdersPlacedCounter()
	ordersApprovedCounter := adapters.NewOrdersApprovedCounter()
	ordersRejectedCounter := adapters.NewOrdersRejectedCounter()

	application := Application{
		Commands: Commands{
			CreateOrder:          commands.NewCreateOrderHandler(orderRepo, restaurantRepo, svc.Logger),
			ApproveOrder:         commands.NewApproveOrderHandler(orderRepo, ordersApprovedCounter),
			RejectOrder:          commands.NewRejectOrderHandler(orderRepo, ordersRejectedCounter),
			BeginCancelOrder:     commands.NewBeginCancelOrderHandler(orderRepo),
			UndoCancelOrder:      commands.NewUndoCancelOrderHandler(orderRepo),
			ConfirmCancelOrder:   commands.NewConfirmCancelOrderHandler(orderRepo),
			BeginReviseOrder:     commands.NewBeginReviseOrderHandler(orderRepo),
			UndoReviseOrder:      commands.NewUndoReviseOrderHandler(orderRepo),
			ConfirmReviseOrder:   commands.NewConfirmReviseOrderHandler(orderRepo),
			StartCreateOrderSaga: commands.NewStartCreateOrderSagaHandler(createOrderSaga, ordersPlacedCounter),
			StartCancelOrderSaga: commands.NewStartCancelOrderSagaHandler(orderRepo, cancelOrderSaga),
			StartReviseOrderSaga: commands.NewStartReviseOrderSagaHandler(orderRepo, reviseOrderSaga),
			CreateRestaurant:     commands.NewCreateRestaurantHandler(restaurantRepo),
			ReviseRestaurantMenu: commands.NewReviseRestaurantMenuHandler(restaurantRepo),
		},
		Queries: Queries{
			GetOrder: queries.NewGetOrderHandler(orderRepo),
		},
	}

	cmdHandlers := NewCommandHandlers(application)
	svc.Subscriber.Subscribe(orderapi.OrderServiceCommandChannel, saga.NewCommandDispatcher(svc.Publisher).
		Handle(orderapi.RejectOrder{}, cmdHandlers.RejectOrder).
		Handle(orderapi.ApproveOrder{}, cmdHandlers.ApproveOrder).
		Handle(orderapi.BeginCancelOrder{}, cmdHandlers.BeginCancel).
		Handle(orderapi.UndoCancelOrder{}, cmdHandlers.UndoCancel).
		Handle(orderapi.ConfirmCancelOrder{}, cmdHandlers.ConfirmCancel).
		Handle(orderapi.BeginReviseOrder{}, cmdHandlers.BeginRevise).
		Handle(orderapi.UndoReviseOrder{}, cmdHandlers.UndoRevise).
		Handle(orderapi.ConfirmReviseOrder{}, cmdHandlers.ConfirmRevise))

	restaurantEventHandlers := newRestaurantEventHandlers(application)
	svc.Subscriber.Subscribe(restaurantapi.RestaurantAggregateChannel, msg.NewEntityEventDispatcher().
		Handle(restaurantapi.RestaurantCreated{}, restaurantEventHandlers.RestaurantCreated).
		Handle(restaurantapi.RestaurantMenuRevised{}, restaurantEventHandlers.RestaurantMenuRevised))

	orderEventHandlers := newOrderEventHandlers(application)
	svc.Subscriber.Subscribe(orderapi.OrderAggregateChannel, msg.NewEntityEventDispatcher().
		Handle(orderapi.OrderCreated{}, orderEventHandlers.OrderCreated))

	orderpb.RegisterOrderServiceServer(svc.RpcServer, newRpcHandlers(application))

	return nil
}
