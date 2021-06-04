package kitcmod

import (
	edatpgx "github.com/stackus/edat-pgx"
	"github.com/stackus/edat/outbox"

	"github.com/stackus/ftgogo/kitchen/internal/adapters"
	"github.com/stackus/ftgogo/kitchen/internal/application"
	"github.com/stackus/ftgogo/kitchen/internal/application/commands"
	"github.com/stackus/ftgogo/kitchen/internal/application/queries"
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
	messageStore := edatpgx.NewMessageStore(svc.PgConn, edatpgx.WithMessageStoreTableName("kitchen.messages"))

	// Driven
	ticketRepo := adapters.NewTicketRepositoryPublisherMiddleware(
		adapters.NewTicketAggregateRootRepository(aggregateStore),
		adapters.NewTicketEntityEventPublisher(svc.Publisher),
	)
	adapters.RestaurantsTableName = "kitchen.restaurants"
	restaurantRepo := adapters.NewRestaurantPostgresRepository(svc.PgConn)

	app := application.Service{
		Commands: application.Commands{
			CreateTicket:         commands.NewCreateTicketHandler(ticketRepo),
			ConfirmCreateTicket:  commands.NewConfirmCreateTicketHandler(ticketRepo),
			CancelCreateTicket:   commands.NewCancelCreateTicketHandler(ticketRepo),
			BeginCancelTicket:    commands.NewBeginCancelTicketHandler(ticketRepo),
			ConfirmCancelTicket:  commands.NewConfirmCancelTicketHandler(ticketRepo),
			UndoCancelTicket:     commands.NewUndoCancelTicketHandler(ticketRepo),
			BeginReviseTicket:    commands.NewBeginReviseTicketHandler(ticketRepo),
			ConfirmReviseTicket:  commands.NewConfirmReviseTicketHandler(ticketRepo),
			UndoReviseTicket:     commands.NewUndoReviseTicketHandler(ticketRepo),
			AcceptTicket:         commands.NewAcceptTicketHandler(ticketRepo),
			CreateRestaurant:     commands.NewCreateRestaurantHandler(restaurantRepo),
			ReviseRestaurantMenu: commands.NewReviseRestaurantMenuHandler(restaurantRepo),
		},
		Queries: application.Queries{
			GetRestaurant: queries.NewGetRestaurantHandler(restaurantRepo),
		},
	}

	// Drivers
	handlers.NewCommandHandlers(app).Mount(svc.Subscriber, svc.Publisher)
	handlers.NewRestaurantEventHandlers(app).Mount(svc.Subscriber)
	handlers.NewRpcHandlers(app).Mount(svc.RpcServer)
	svc.Processors = append(svc.Processors, outbox.NewPollingProcessor(messageStore, svc.CDCPublisher))

	return nil
}
