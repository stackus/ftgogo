package main

import (
	"github.com/stackus/edat/msg"

	"github.com/stackus/ftgogo/delivery/internal/adapters"
	"github.com/stackus/ftgogo/delivery/internal/application/commands"
	"github.com/stackus/ftgogo/delivery/internal/application/queries"
	"github.com/stackus/ftgogo/serviceapis"
	"github.com/stackus/ftgogo/serviceapis/deliveryapi/pb"
	"github.com/stackus/ftgogo/serviceapis/kitchenapi"
	"github.com/stackus/ftgogo/serviceapis/orderapi"
	"github.com/stackus/ftgogo/serviceapis/restaurantapi"
	"shared-go/applications"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateDelivery         commands.CreateDeliveryHandler
	CreateRestaurant       commands.CreateRestaurantHandler
	SetCourierAvailability commands.SetCourierAvailabilityHandler
	ScheduleDelivery       commands.ScheduleDeliveryHandler
	CancelDelivery         commands.CancelDeliveryHandler
}

type Queries struct {
	GetDeliveryStatus queries.GetDeliveryStatusHandler
}

func main() {
	svc := applications.NewService(initService)
	if err := svc.Execute(); err != nil {
		panic(err)
	}
}

func initService(svc *applications.Service) error {
	serviceapis.RegisterTypes()

	courierRepo := adapters.NewCourierPostgresRepository(svc.PgConn)
	deliveryRepo := adapters.NewDeliveryPostgresRepository(svc.PgConn)
	restaurantRepo := adapters.NewRestaurantPostgresRepository(svc.PgConn)

	application := Application{
		Commands: Commands{
			CreateDelivery:         commands.NewCreateDeliveryHandler(deliveryRepo, restaurantRepo),
			CreateRestaurant:       commands.NewCreateRestaurantHandler(restaurantRepo),
			SetCourierAvailability: commands.NewSetCourierAvailabilityHandler(courierRepo),
			ScheduleDelivery:       commands.NewScheduleDeliveryHandler(deliveryRepo, courierRepo),
			CancelDelivery:         commands.NewCancelDeliveryHandler(deliveryRepo, courierRepo),
		},
		Queries: Queries{
			GetDeliveryStatus: queries.NewGetDeliveryStatusHandler(deliveryRepo, courierRepo),
		},
	}

	restaurantEventHandlers := newRestaurantEventHandlers(application)
	svc.Subscriber.Subscribe(restaurantapi.RestaurantAggregateChannel, msg.NewEntityEventDispatcher().
		Handle(restaurantapi.RestaurantCreated{}, restaurantEventHandlers.RestaurantCreated))

	orderEventHandlers := newOrderEventHandlers(application)
	svc.Subscriber.Subscribe(orderapi.OrderAggregateChannel, msg.NewEntityEventDispatcher().
		Handle(orderapi.OrderCreated{}, orderEventHandlers.OrderCreated))

	ticketEventHandlers := newTicketEventHandlers(application)
	svc.Subscriber.Subscribe(kitchenapi.TicketAggregateChannel, msg.NewEntityEventDispatcher().
		Handle(kitchenapi.TicketAccepted{}, ticketEventHandlers.TicketAccepted).
		Handle(kitchenapi.TicketCancelled{}, ticketEventHandlers.TicketCancelled))

	deliverypb.RegisterDeliveryServiceServer(svc.RpcServer, newRpcHandlers(application))

	return nil
}
