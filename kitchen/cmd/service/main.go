package main

import (
	"github.com/stackus/ftgogo/kitchen/internal/adapters"
	"github.com/stackus/ftgogo/kitchen/internal/application"
	"github.com/stackus/ftgogo/kitchen/internal/application/commands"
	"github.com/stackus/ftgogo/kitchen/internal/application/queries"
	"github.com/stackus/ftgogo/kitchen/internal/domain"
	"github.com/stackus/ftgogo/kitchen/internal/handlers"
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
	ticketRepo := adapters.NewTicketRepositoryPublisherMiddleware(
		adapters.NewTicketAggregateRepository(svc.AggregateStore),
		adapters.NewTicketEntityEventPublisher(svc.Publisher),
	)
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

	return nil
}
