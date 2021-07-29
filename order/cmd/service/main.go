package main

import (
	"github.com/stackus/ftgogo/order/internal/adapters"
	"github.com/stackus/ftgogo/order/internal/application"
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
		adapters.NewOrderAggregateRepository(svc.AggregateStore),
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
	ordersPlacedCounter := adapters.NewPrometheusCounter("orders_placed")
	ordersApprovedCounter := adapters.NewPrometheusCounter("orders_approved")
	ordersRejectedCounter := adapters.NewPrometheusCounter("orders_rejected")

	app := application.NewServiceApplication(
		orderRepo, restaurantRepo,
		createOrderSaga, cancelOrderSaga, reviseOrderSaga,
		ordersPlacedCounter, ordersApprovedCounter, ordersRejectedCounter,
	)

	// Drivers
	handlers.NewCommandHandlers(app).Mount(svc.Subscriber, svc.Publisher)
	handlers.NewRestaurantEventHandlers(app).Mount(svc.Subscriber)
	handlers.NewOrderEventHandlers(app).Mount(svc.Subscriber)
	handlers.NewRpcHandlers(app).Mount(svc.RpcServer)

	return nil
}
