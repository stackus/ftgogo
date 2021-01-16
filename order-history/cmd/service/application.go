package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"github.com/stackus/edat/msg"

	"github.com/stackus/ftgogo/order-history/internal/adapters"
	"github.com/stackus/ftgogo/order-history/internal/application/commands"
	"github.com/stackus/ftgogo/order-history/internal/application/queries"
	"serviceapis"
	"serviceapis/orderapi"
	"shared-go/applications"
	"shared-go/web"
)

// To regenerate the web server api use the following generate command
//go:generate oapi-codegen -config oapi-codegen.cfg.yaml openapi.yaml

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateOrderHistory commands.CreateOrderHistoryHandler
	UpdateOrderStatus  commands.UpdateOrderStatusHandler
}

type Queries struct {
	GetConsumerOrderHistory queries.GetConsumerOrderHistoryHandler
	GetOrderHistory         queries.GetOrderHistoryHandler
}

func initApplication(svc *applications.Service) error {
	serviceapis.RegisterTypes()

	orderHistoryRepo := adapters.NewOrderHistoryPostgresRepository(svc.PgConn)

	application := Application{
		Commands: Commands{
			CreateOrderHistory: commands.NewCreateOrderHistoryHandler(orderHistoryRepo),
			UpdateOrderStatus:  commands.NewUpdateOrderStatusHandler(orderHistoryRepo),
		},
		Queries: Queries{
			GetConsumerOrderHistory: queries.NewGetConsumerOrderHistoryHandler(orderHistoryRepo),
			GetOrderHistory:         queries.NewGetOrderHistoryHandler(orderHistoryRepo),
		},
	}

	orderEventHandlers := NewOrderEventHandlers(application)
	svc.Subscriber.Subscribe(orderapi.OrderAggregateChannel, msg.NewEntityEventDispatcher().
		Handle(orderapi.OrderCreated{}, orderEventHandlers.OrderCreated).
		Handle(orderapi.OrderApproved{}, orderEventHandlers.OrderApproved).
		Handle(orderapi.OrderCancelled{}, orderEventHandlers.OrderCancelled).
		Handle(orderapi.OrderRejected{}, orderEventHandlers.OrderRejected))

	svc.WebServer.Mount(svc.Cfg.Web.ApiPath, func(r chi.Router) http.Handler {
		return HandlerFromMux(NewWebHandlers(application), r)
	})

	return nil
}

type OrderEventHandlers struct{ app Application }

func NewOrderEventHandlers(app Application) OrderEventHandlers { return OrderEventHandlers{app: app} }

func (h OrderEventHandlers) OrderCreated(ctx context.Context, evtMsg msg.EntityEvent) error {
	evt := evtMsg.Event().(*orderapi.OrderCreated)

	return h.app.Commands.CreateOrderHistory.Handle(ctx, commands.CreateOrderHistory{
		OrderID:        evtMsg.EntityID(),
		ConsumerID:     evt.ConsumerID,
		RestaurantID:   evt.RestaurantID,
		RestaurantName: evt.RestaurantName,
		LineItems:      evt.LineItems,
		OrderTotal:     evt.OrderTotal,
	})
}

func (h OrderEventHandlers) OrderApproved(ctx context.Context, evtMsg msg.EntityEvent) error {
	return h.app.Commands.UpdateOrderStatus.Handle(ctx, commands.UpdateOrderStatus{
		OrderID: evtMsg.EntityID(),
		Status:  orderapi.Approved,
	})
}

func (h OrderEventHandlers) OrderCancelled(ctx context.Context, evtMsg msg.EntityEvent) error {
	return h.app.Commands.UpdateOrderStatus.Handle(ctx, commands.UpdateOrderStatus{
		OrderID: evtMsg.EntityID(),
		Status:  orderapi.Cancelled,
	})
}

func (h OrderEventHandlers) OrderRejected(ctx context.Context, evtMsg msg.EntityEvent) error {
	return h.app.Commands.UpdateOrderStatus.Handle(ctx, commands.UpdateOrderStatus{
		OrderID: evtMsg.EntityID(),
		Status:  orderapi.Rejected,
	})
}

type WebHandlers struct {
	app Application
}

func NewWebHandlers(app Application) WebHandlers {
	return WebHandlers{app: app}
}

func (h WebHandlers) GetConsumerOrderHistory(w http.ResponseWriter, r *http.Request, params GetConsumerOrderHistoryParams) {
	if response, err := h.app.Queries.GetConsumerOrderHistory.Handle(r.Context(), queries.GetConsumerOrderHistory{
		ConsumerID: params.ConsumerID,
		Filter:     params.Filter,
		Next:       params.Next,
		Limit:      params.Limit,
	}); err != nil {
		render.Render(w, r, web.NewErrorResponse(err))
		return
	} else {
		render.Respond(w, r, response)
	}
}

func (h WebHandlers) GetOrderHistory(w http.ResponseWriter, r *http.Request, orderID OrderID) {
	if response, err := h.app.Queries.GetOrderHistory.Handle(r.Context(), queries.GetOrderHistory{
		OrderID: string(orderID),
	}); err != nil {
		render.Render(w, r, web.NewErrorResponse(err))
		return
	} else {
		render.Respond(w, r, response)
	}
}
