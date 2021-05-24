package main

import (
	"github.com/stackus/edat/msg"
	"github.com/stackus/edat/saga"

	"github.com/stackus/ftgogo/kitchen/internal/adapters"
	"github.com/stackus/ftgogo/kitchen/internal/application/commands"
	"github.com/stackus/ftgogo/kitchen/internal/application/queries"
	"github.com/stackus/ftgogo/kitchen/internal/domain"
	"github.com/stackus/ftgogo/serviceapis"
	"github.com/stackus/ftgogo/serviceapis/kitchenapi"
	"github.com/stackus/ftgogo/serviceapis/kitchenapi/pb"
	"github.com/stackus/ftgogo/serviceapis/restaurantapi"
	"shared-go/applications"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateTicket         commands.CreateTicketHandler
	ConfirmCreateTicket  commands.ConfirmCreateTicketHandler
	CancelCreateTicket   commands.CancelCreateTicketHandler
	BeginCancelTicket    commands.BeginCancelTicketHandler
	ConfirmCancelTicket  commands.ConfirmCancelTicketHandler
	UndoCancelTicket     commands.UndoCancelTicketHandler
	BeginReviseTicket    commands.BeginReviseTicketHandler
	ConfirmReviseTicket  commands.ConfirmReviseTicketHandler
	UndoReviseTicket     commands.UndoReviseTicketHandler
	AcceptTicket         commands.AcceptTicketHandler
	CreateRestaurant     commands.CreateRestaurantHandler
	ReviseRestaurantMenu commands.ReviseRestaurantMenuHandler
}

type Queries struct {
	GetRestaurant queries.GetRestaurantHandler
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

	ticketRepo := adapters.NewTicketRepositoryPublisherMiddleware(
		adapters.NewTicketRepository(svc.AggregateStore),
		adapters.NewTicketPublisher(svc.Publisher),
	)
	restaurantRepo := adapters.NewRestaurantPostgresRepository(svc.PgConn)

	application := Application{
		Commands: Commands{
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
		Queries: Queries{
			GetRestaurant: queries.NewGetRestaurantHandler(restaurantRepo),
		},
	}

	cmdHandlers := NewCommandHandlers(application)
	svc.Subscriber.Subscribe(kitchenapi.KitchenServiceCommandChannel, saga.NewCommandDispatcher(svc.Publisher).
		Handle(kitchenapi.CreateTicket{}, cmdHandlers.CreateTicket).
		Handle(kitchenapi.ConfirmCreateTicket{}, cmdHandlers.ConfirmCreateTicket).
		Handle(kitchenapi.CancelCreateTicket{}, cmdHandlers.CancelCreateTicket).
		Handle(kitchenapi.BeginCancelTicket{}, cmdHandlers.BeginCancelTicket).
		Handle(kitchenapi.ConfirmCancelTicket{}, cmdHandlers.ConfirmCancelTicket).
		Handle(kitchenapi.UndoCancelTicket{}, cmdHandlers.UndoCancelTicket).
		Handle(kitchenapi.BeginReviseTicket{}, cmdHandlers.BeginReviseTicket).
		Handle(kitchenapi.ConfirmReviseTicket{}, cmdHandlers.ConfirmReviseTicket).
		Handle(kitchenapi.UndoReviseTicket{}, cmdHandlers.UndoReviseTicket))

	restaurantEventHandlers := newRestaurantEventHandlers(application)
	svc.Subscriber.Subscribe(restaurantapi.RestaurantAggregateChannel, msg.NewEntityEventDispatcher().
		Handle(restaurantapi.RestaurantCreated{}, restaurantEventHandlers.RestaurantCreated).
		Handle(restaurantapi.RestaurantMenuRevised{}, restaurantEventHandlers.RestaurantMenuRevised))

	kitchenpb.RegisterKitchenServiceServer(svc.RpcServer, newRpcHandlers(application))

	return nil
}
