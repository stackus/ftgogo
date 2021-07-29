package ordmod

import (
	edatpgx "github.com/stackus/edat-pgx"
	"github.com/stackus/edat/msg"
	"github.com/stackus/edat/outbox"

	"github.com/stackus/ftgogo/order/internal/adapters"
	"github.com/stackus/ftgogo/order/internal/application"
	"github.com/stackus/ftgogo/order/internal/domain"
	"github.com/stackus/ftgogo/order/internal/handlers"
	"shared-go/applications"
)

func Setup(svc *applications.Monolith) error {
	domain.RegisterTypes()

	// Infrastructure
	aggregateStore := edatpgx.NewSnapshotStore(
		svc.PgConn,
		edatpgx.WithSnapshotStoreTableName("orders.snapshots"),
	)(edatpgx.NewEventStore(
		svc.PgConn,
		edatpgx.WithEventStoreTableName("orders.events"),
	))
	sagaInstanceStore := edatpgx.NewSagaInstanceStore(svc.PgConn, edatpgx.WithSagaInstanceStoreTableName("orders.saga_instances"))
	messageStore := edatpgx.NewMessageStore(svc.CDCPgConn, edatpgx.WithMessageStoreTableName("orders.messages"))
	publisher := msg.NewPublisher(messageStore)
	svc.Publishers = append(svc.Publishers, publisher)
	svc.Processors = append(svc.Processors, outbox.NewPollingProcessor(messageStore, svc.CDCPublisher))

	// Driven
	orderRepo := adapters.NewOrderRepositoryPublisherMiddleware(
		adapters.NewOrderAggregateRepository(aggregateStore),
		adapters.NewOrderEntityEventPublisher(publisher),
	)
	adapters.RestaurantsTableName = "orders.restaurants"
	restaurantRepo := adapters.NewRestaurantPostgresRepository(svc.PgConn)

	createOrderSaga := adapters.NewCreateOrderOrchestrationSaga(sagaInstanceStore, publisher)
	svc.Subscriber.Subscribe(createOrderSaga.ReplyChannel(), createOrderSaga)

	cancelOrderSaga := adapters.NewCancelOrderOrchestrationSaga(sagaInstanceStore, publisher)
	svc.Subscriber.Subscribe(cancelOrderSaga.ReplyChannel(), cancelOrderSaga)

	reviseOrderSaga := adapters.NewReviseOrderOrchestrationSaga(sagaInstanceStore, publisher)
	svc.Subscriber.Subscribe(reviseOrderSaga.ReplyChannel(), reviseOrderSaga)

	ordersPlacedCounter := adapters.NewPrometheusCounter("orders_placed")
	ordersApprovedCounter := adapters.NewPrometheusCounter("orders_approved")
	ordersRejectedCounter := adapters.NewPrometheusCounter("orders_rejected")

	app := application.NewServiceApplication(
		orderRepo, restaurantRepo,
		createOrderSaga, cancelOrderSaga, reviseOrderSaga,
		ordersPlacedCounter, ordersApprovedCounter, ordersRejectedCounter,
	)

	// Drivers
	handlers.NewCommandHandlers(app).Mount(svc.Subscriber, publisher)
	handlers.NewRestaurantEventHandlers(app).Mount(svc.Subscriber)
	handlers.NewOrderEventHandlers(app).Mount(svc.Subscriber)
	handlers.NewRpcHandlers(app).Mount(svc.RpcServer)

	return nil
}
