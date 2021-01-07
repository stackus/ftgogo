package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"github.com/stackus/edat/msg"
	"github.com/stackus/ftgogo/delivery/internal/adapters"
	"github.com/stackus/ftgogo/delivery/internal/application/commands"
	"github.com/stackus/ftgogo/delivery/internal/application/queries"
	"serviceapis/kitchenapi"
	"serviceapis/orderapi"
	"serviceapis/restaurantapi"
	"shared-go/applications"
	"shared-go/web"

	"serviceapis"
)

// To regenerate the web server api use the following generate command
//go:generate oapi-codegen -config oapi-codegen.cfg.yaml openapi.yaml

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

func initApplication(svc *applications.Service) error {
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

	restaurantEventHandlers := NewRestaurantEventHandlers(application)
	svc.Subscriber.Subscribe(restaurantapi.RestaurantAggregateChannel, msg.NewEntityEventDispatcher().
		Handle(restaurantapi.RestaurantCreated{}, restaurantEventHandlers.RestaurantCreated))

	orderEventHandlers := NewOrderEventHandlers(application)
	svc.Subscriber.Subscribe(orderapi.OrderAggregateChannel, msg.NewEntityEventDispatcher().
		Handle(orderapi.OrderCreated{}, orderEventHandlers.OrderCreated))

	ticketEventHandlers := NewTicketEventHandlers(application)
	svc.Subscriber.Subscribe(kitchenapi.TicketAggregateChannel, msg.NewEntityEventDispatcher().
		Handle(kitchenapi.TicketAccepted{}, ticketEventHandlers.TicketAccepted).
		Handle(kitchenapi.TicketCancelled{}, ticketEventHandlers.TicketCancelled))

	// TODO refactor so a string isn't used here
	svc.WebServer.Mount("/api", func(r chi.Router) http.Handler {
		return HandlerFromMux(NewWebHandlers(application), r)
	})

	return nil
}

type WebHandlers struct{ app Application }

func NewWebHandlers(app Application) WebHandlers { return WebHandlers{app: app} }

func (h WebHandlers) SetCourierAvailability(w http.ResponseWriter, r *http.Request, courierID CourierID) {
	cid := string(courierID)

	request := SetCourierAvailabilityJSONRequestBody{}

	if err := render.Decode(r, &request); err != nil {
		render.Render(w, r, web.NewErrorResponse(err))
		return
	}

	err := h.app.Commands.SetCourierAvailability.Handle(r.Context(), commands.SetCourierAvailability{
		CourierID: cid,
		Available: request.Available,
	})
	if err != nil {
		render.Render(w, r, web.NewErrorResponse(err))
		return
	}

	render.Respond(w, r, CourierAvailableResponse{Available: request.Available})
}

func (h WebHandlers) GetDeliveryStatus(w http.ResponseWriter, r *http.Request, deliveryID DeliveryID) {
	did := string(deliveryID)

	status, err := h.app.Queries.GetDeliveryStatus.Handle(r.Context(), queries.GetDeliveryStatus{DeliveryID: did})
	if err != nil {
		render.Render(w, r, web.NewErrorResponse(err))
		return
	}

	render.Respond(w, r, status)
}

type RestaurantEventHandlers struct{ app Application }

func NewRestaurantEventHandlers(app Application) RestaurantEventHandlers {
	return RestaurantEventHandlers{app: app}
}

func (h RestaurantEventHandlers) RestaurantCreated(ctx context.Context, evtMsg msg.EntityEvent) error {
	evt := evtMsg.Event().(*restaurantapi.RestaurantCreated)

	return h.app.Commands.CreateRestaurant.Handle(ctx, commands.CreateRestaurant{
		RestaurantID: evtMsg.EntityID(),
		Name:         evt.Name,
		Address:      evt.Address,
	})
}

type OrderEventHandlers struct{ app Application }

func NewOrderEventHandlers(app Application) OrderEventHandlers { return OrderEventHandlers{app: app} }

func (h OrderEventHandlers) OrderCreated(ctx context.Context, evtMsg msg.EntityEvent) error {
	evt := evtMsg.Event().(*orderapi.OrderCreated)

	return h.app.Commands.CreateDelivery.Handle(ctx, commands.CreateDelivery{
		OrderID:         evtMsg.EntityID(),
		RestaurantID:    evt.RestaurantID,
		DeliveryAddress: evt.DeliverTo,
	})
}

type TicketEventHandlers struct{ app Application }

func NewTicketEventHandlers(app Application) TicketEventHandlers {
	return TicketEventHandlers{app: app}
}

func (h TicketEventHandlers) TicketAccepted(ctx context.Context, evtMsg msg.EntityEvent) error {
	evt := evtMsg.Event().(*kitchenapi.TicketAccepted)

	return h.app.Commands.ScheduleDelivery.Handle(ctx, commands.ScheduleDelivery{
		OrderID: evt.OrderID,
		ReadyBy: evt.ReadyBy,
	})
}

func (h TicketEventHandlers) TicketCancelled(ctx context.Context, evtMsg msg.EntityEvent) error {
	evt := evtMsg.Event().(*kitchenapi.TicketCancelled)

	return h.app.Commands.CancelDelivery.Handle(ctx, commands.CancelDelivery{OrderID: evt.OrderID})
}
