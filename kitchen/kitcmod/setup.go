package kitcmod

import (
	edatpgx "github.com/stackus/edat-pgx"
	"github.com/stackus/edat/msg"
	"github.com/stackus/edat/outbox"

	"github.com/stackus/ftgogo/kitchen/internal/adapters"
	"github.com/stackus/ftgogo/kitchen/internal/application"
	"github.com/stackus/ftgogo/kitchen/internal/domain"
	"github.com/stackus/ftgogo/kitchen/internal/handlers"
	"shared-go/applications"
)

func Setup(svc *applications.Monolith) error {
	domain.RegisterTypes()

	// Infrastructure
	aggregateStore := edatpgx.NewSnapshotStore(
		svc.PgConn,
		edatpgx.WithSnapshotStoreTableName("kitchen.snapshots"),
	)(edatpgx.NewEventStore(
		svc.PgConn,
		edatpgx.WithEventStoreTableName("kitchen.events"),
	))
	messageStore := edatpgx.NewMessageStore(svc.CDCPgConn, edatpgx.WithMessageStoreTableName("kitchen.messages"))
	publisher := msg.NewPublisher(messageStore)
	svc.Publishers = append(svc.Publishers, publisher)
	svc.Processors = append(svc.Processors, outbox.NewPollingProcessor(messageStore, svc.CDCPublisher))

	// Driven
	ticketRepo := adapters.NewTicketRepositoryPublisherMiddleware(
		adapters.NewTicketAggregateRepository(aggregateStore),
		adapters.NewTicketEntityEventPublisher(publisher),
	)
	adapters.RestaurantsTableName = "kitchen.restaurants"
	restaurantRepo := adapters.NewRestaurantPostgresRepository(svc.PgConn)

	app := application.NewServiceApplication(ticketRepo, restaurantRepo)

	// Drivers
	handlers.NewCommandHandlers(app).Mount(svc.Subscriber, publisher)
	handlers.NewRestaurantEventHandlers(app).Mount(svc.Subscriber)
	handlers.NewRpcHandlers(app).Mount(svc.RpcServer)

	return nil
}
