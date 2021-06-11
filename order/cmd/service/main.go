package main

import (
	"github.com/stackus/ftgogo/order/internal/adapters"
	"github.com/stackus/ftgogo/order/internal/application"
	"github.com/stackus/ftgogo/order/internal/application/commands"
	"github.com/stackus/ftgogo/order/internal/application/queries"
	"github.com/stackus/ftgogo/order/internal/domain"
	"github.com/stackus/ftgogo/order/internal/handlers"
	"github.com/stackus/ftgogo/serviceapis"
	"shared-go/applications"
)

func main() {
	svc := applications.NewService(initService)
	if err := svc.Execute(); err != nil {
		panic(err)
	}
}

func initService(svc *applications.Service) error {
	serviceapis.RegisterTypes()
	domain.RegisterTypes()

	// Driven
	orderRepo := adapters.NewOrderRepositoryPublisherMiddleware(
		adapters.NewOrderAggregateRootRepository(svc.AggregateStore),
		adapters.NewOrderEntityEventPublisher(svc.Publisher),
	)
	restaurantRepo := adapters.NewRestaurantPostgresRepository(svc.PgConn)

	createOrderSaga := adapters.NewCreateOrderOrchestrationSaga(svc.SagaInstanceStore, svc.Publisher)
	svc.Subscriber.Subscribe(createOrderSaga.ReplyChannel(), createOrderSaga)

	cancelOrderSaga := adapters.NewCancelOrderOrchestrationSaga(svc.SagaInstanceStore, svc.Publisher)
	svc.Subscriber.Subscribe(cancelOrderSaga.ReplyChannel(), cancelOrderSaga)

	reviseOrderSaga := adapters.NewReviseOrderOrchestrationSaga(svc.SagaInstanceStore, svc.Publisher)
	svc.Subscriber.Subscribe(reviseOrderSaga.ReplyChannel(), reviseOrderSaga)

	// Counters
	ordersPlacedCounter := adapters.NewOrdersPlacedCounter()
	ordersApprovedCounter := adapters.NewOrdersApprovedCounter()
	ordersRejectedCounter := adapters.NewOrdersRejectedCounter()

	app := application.Service{
		Commands: application.Commands{
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
		Queries: application.Queries{
			GetOrder:      queries.NewGetOrderHandler(orderRepo),
			GetRestaurant: queries.NewGetRestaurantHandler(restaurantRepo),
		},
	}

	// Drivers
	handlers.NewCommandHandlers(app).Mount(svc.Subscriber, svc.Publisher)
	handlers.NewRestaurantEventHandlers(app).Mount(svc.Subscriber)
	handlers.NewOrderEventHandlers(app).Mount(svc.Subscriber)
	handlers.NewRpcHandlers(app).Mount(svc.RpcServer)

	return nil
}
