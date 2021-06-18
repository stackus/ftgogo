package ordmod

import (
	edatpgx "github.com/stackus/edat-pgx"
	"github.com/stackus/edat/msg"
	"github.com/stackus/edat/outbox"

	"github.com/stackus/ftgogo/order/internal/adapters"
	"github.com/stackus/ftgogo/order/internal/application"
	"github.com/stackus/ftgogo/order/internal/application/commands"
	"github.com/stackus/ftgogo/order/internal/application/queries"
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
	handlers.NewCommandHandlers(app).Mount(svc.Subscriber, publisher)
	handlers.NewRestaurantEventHandlers(app).Mount(svc.Subscriber)
	handlers.NewOrderEventHandlers(app).Mount(svc.Subscriber)
	handlers.NewRpcHandlers(app).Mount(svc.RpcServer)

	return nil
}
