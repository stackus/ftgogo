package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/stackus/edat/msg"
	"github.com/stackus/edat/saga"

	"github.com/stackus/ftgogo/order/internal/adapters"
	"github.com/stackus/ftgogo/order/internal/application/commands"
	"github.com/stackus/ftgogo/order/internal/application/queries"
	"github.com/stackus/ftgogo/order/internal/domain"
	"serviceapis"
	"serviceapis/orderapi"
	"serviceapis/restaurantapi"
	"shared-go/applications"
	"shared-go/web"
)

// To regenerate the web server api use the following generate command
//go:generate oapi-codegen -config oapi-codegen.cfg.yaml openapi.yaml

type Application struct {
	Commands Commands
	Queries  Queries
	Metrics  Metrics
}

type Commands struct {
	CreateOrder          commands.CreateOrderHandler
	ApproveOrder         commands.ApproveOrderHandler
	RejectOrder          commands.RejectOrderHandler
	BeginCancelOrder     commands.BeginCancelOrderHandler
	UndoCancelOrder      commands.UndoCancelOrderHandler
	ConfirmCancelOrder   commands.ConfirmCancelOrderHandler
	BeginReviseOrder     commands.BeginReviseOrderHandler
	UndoReviseOrder      commands.UndoReviseOrderHandler
	ConfirmReviseOrder   commands.ConfirmReviseOrderHandler
	StartCreateOrderSaga commands.StartCreateOrderSagaHandler
	StartCancelOrderSaga commands.StartCancelOrderSagaHandler
	StartReviseOrderSaga commands.StartReviseOrderSagaHandler
	CreateRestaurant     commands.CreateRestaurantHandler
	ReviseRestaurantMenu commands.ReviseRestaurantMenuHandler
}

type Queries struct {
	GetOrder      queries.GetOrderHandler
	GetRestaurant queries.GetRestaurantHandler
}

type Metrics struct {
	OrdersPlaced   prometheus.Counter
	OrdersApproved prometheus.Counter
	OrdersRejected prometheus.Counter
}

func initApplication(svc *applications.Service) error {
	serviceapis.RegisterTypes()
	domain.RegisterTypes()

	orderRepo := adapters.NewOrderRepository(svc.AggregateStore)
	orderPublisher := adapters.NewOrderPublisher(svc.Publisher)
	restaurantRepo := adapters.NewRestaurantPostgresRepository(svc.PgConn)

	createOrderSaga := adapters.NewCreateOrderSaga(svc.SagaInstanceStore, svc.Publisher)
	svc.Subscriber.Subscribe(createOrderSaga.ReplyChannel(), createOrderSaga)

	cancelOrderSaga := adapters.NewCancelOrderSaga(svc.SagaInstanceStore, svc.Publisher)
	svc.Subscriber.Subscribe(cancelOrderSaga.ReplyChannel(), cancelOrderSaga)

	reviseOrderSaga := adapters.NewReviseOrderSaga(svc.SagaInstanceStore, svc.Publisher)
	svc.Subscriber.Subscribe(reviseOrderSaga.ReplyChannel(), reviseOrderSaga)

	application := Application{
		Commands: Commands{
			CreateOrder:          commands.NewCreateOrderHandler(orderRepo, restaurantRepo, orderPublisher, svc.Logger),
			ApproveOrder:         commands.NewApproveOrderHandler(orderRepo, orderPublisher),
			RejectOrder:          commands.NewRejectOrderHandler(orderRepo, orderPublisher),
			BeginCancelOrder:     commands.NewBeginCancelOrderHandler(orderRepo),
			UndoCancelOrder:      commands.NewUndoCancelOrderHandler(orderRepo),
			ConfirmCancelOrder:   commands.NewConfirmCancelOrderHandler(orderRepo, orderPublisher),
			BeginReviseOrder:     commands.NewBeginReviseOrderHandler(orderRepo, orderPublisher),
			UndoReviseOrder:      commands.NewUndoReviseOrderHandler(orderRepo),
			ConfirmReviseOrder:   commands.NewConfirmReviseOrderHandler(orderRepo, orderPublisher),
			StartCreateOrderSaga: commands.NewStartCreateOrderSagaHandler(createOrderSaga),
			StartCancelOrderSaga: commands.NewStartCancelOrderSagaHandler(orderRepo, cancelOrderSaga),
			StartReviseOrderSaga: commands.NewStartReviseOrderSagaHandler(orderRepo, reviseOrderSaga),
			CreateRestaurant:     commands.NewCreateRestaurantHandler(restaurantRepo),
			ReviseRestaurantMenu: commands.NewReviseRestaurantMenuHandler(restaurantRepo),
		},
		Queries: Queries{
			GetOrder:      queries.NewGetOrderHandler(orderRepo),
			GetRestaurant: queries.NewGetRestaurantHandler(restaurantRepo),
		},
		Metrics: Metrics{
			OrdersPlaced:   promauto.NewCounter(prometheus.CounterOpts{Name: "placed_orders"}),
			OrdersApproved: promauto.NewCounter(prometheus.CounterOpts{Name: "approved_orders"}),
			OrdersRejected: promauto.NewCounter(prometheus.CounterOpts{Name: "rejected_orders"}),
		},
	}

	cmdHandlers := NewCommandHandlers(application)
	svc.Subscriber.Subscribe(orderapi.OrderServiceCommandChannel, saga.NewCommandDispatcher(svc.Publisher).
		Handle(orderapi.RejectOrder{}, cmdHandlers.RejectOrder).
		Handle(orderapi.ApproveOrder{}, cmdHandlers.ApproveOrder).
		Handle(orderapi.BeginCancelOrder{}, cmdHandlers.BeginCancel).
		Handle(orderapi.UndoCancelOrder{}, cmdHandlers.UndoCancel).
		Handle(orderapi.ConfirmCancelOrder{}, cmdHandlers.ConfirmCancel).
		Handle(orderapi.BeginReviseOrder{}, cmdHandlers.BeginRevise).
		Handle(orderapi.UndoReviseOrder{}, cmdHandlers.UndoRevise).
		Handle(orderapi.ConfirmReviseOrder{}, cmdHandlers.ConfirmRevise))

	restaurantEventHandlers := NewRestaurantEventHandlers(application)
	svc.Subscriber.Subscribe(restaurantapi.RestaurantAggregateChannel, msg.NewEntityEventDispatcher().
		Handle(restaurantapi.RestaurantCreated{}, restaurantEventHandlers.RestaurantCreated).
		Handle(restaurantapi.RestaurantMenuRevised{}, restaurantEventHandlers.RestaurantMenuRevised))

	orderEventHandlers := NewOrderEventHandlers(application)
	svc.Subscriber.Subscribe(orderapi.OrderAggregateChannel, msg.NewEntityEventDispatcher().
		Handle(orderapi.OrderCreated{}, orderEventHandlers.OrderCreated))

	svc.WebServer.Mount(svc.Cfg.Web.ApiPath, func(r chi.Router) http.Handler {
		return HandlerFromMux(NewWebHandlers(application), r)
	})

	// TODO refactor into an option
	svc.WebServer.Mount("/swagger", func(router chi.Router) http.Handler {
		router.Get("/", func(w http.ResponseWriter, r *http.Request) {
			swagger, err := GetSwagger()
			if err != nil {
				render.Status(r, http.StatusInternalServerError)
				render.PlainText(w, r, fmt.Sprintf("Error rendering Swagger API: %s", err.Error()))
				return
			}
			render.JSON(w, r, swagger)
		})
		return router
	})

	return nil
}

type CommandHandlers struct{ app Application }

func NewCommandHandlers(app Application) CommandHandlers { return CommandHandlers{app: app} }

func (h CommandHandlers) RejectOrder(ctx context.Context, cmdMsg saga.Command) ([]msg.Reply, error) {
	cmd := cmdMsg.Command().(*orderapi.RejectOrder)

	err := h.app.Commands.RejectOrder.Handle(ctx, commands.RejectOrder{OrderID: cmd.OrderID})
	if err != nil {
		return []msg.Reply{msg.WithFailure()}, nil
	}

	h.app.Metrics.OrdersRejected.Inc()

	return []msg.Reply{msg.WithSuccess()}, nil
}

func (h CommandHandlers) ApproveOrder(ctx context.Context, cmdMsg saga.Command) ([]msg.Reply, error) {
	cmd := cmdMsg.Command().(*orderapi.ApproveOrder)

	err := h.app.Commands.ApproveOrder.Handle(ctx, commands.ApproveOrder{
		OrderID:  cmd.OrderID,
		TicketID: cmd.TicketID,
	})
	if err != nil {
		return []msg.Reply{msg.WithFailure()}, nil
	}

	h.app.Metrics.OrdersApproved.Inc()

	return []msg.Reply{msg.WithSuccess()}, nil
}

func (h CommandHandlers) BeginCancel(ctx context.Context, cmdMsg saga.Command) ([]msg.Reply, error) {
	cmd := cmdMsg.Command().(orderapi.BeginCancelOrder)

	err := h.app.Commands.BeginCancelOrder.Handle(ctx, commands.BeginCancelOrder{OrderID: cmd.OrderID})
	if err != nil {
		return []msg.Reply{msg.WithFailure()}, nil
	}

	return []msg.Reply{msg.WithSuccess()}, nil
}

func (h CommandHandlers) UndoCancel(ctx context.Context, cmdMsg saga.Command) ([]msg.Reply, error) {
	cmd := cmdMsg.Command().(orderapi.UndoCancelOrder)

	err := h.app.Commands.UndoCancelOrder.Handle(ctx, commands.UndoCancelOrder{OrderID: cmd.OrderID})
	if err != nil {
		return []msg.Reply{msg.WithFailure()}, nil
	}

	return []msg.Reply{msg.WithSuccess()}, nil
}

func (h CommandHandlers) ConfirmCancel(ctx context.Context, cmdMsg saga.Command) ([]msg.Reply, error) {
	cmd := cmdMsg.Command().(orderapi.ConfirmCancelOrder)

	err := h.app.Commands.ConfirmCancelOrder.Handle(ctx, commands.ConfirmCancelOrder{OrderID: cmd.OrderID})
	if err != nil {
		return []msg.Reply{msg.WithFailure()}, nil
	}

	return []msg.Reply{msg.WithSuccess()}, nil
}

func (h CommandHandlers) BeginRevise(ctx context.Context, cmdMsg saga.Command) ([]msg.Reply, error) {
	cmd := cmdMsg.Command().(orderapi.BeginReviseOrder)

	newTotal, err := h.app.Commands.BeginReviseOrder.Handle(ctx, commands.BeginReviseOrder{
		OrderID:           cmd.OrderID,
		RevisedQuantities: cmd.RevisedQuantities,
	})
	if err != nil {
		return []msg.Reply{msg.WithFailure()}, nil
	}

	return []msg.Reply{msg.WithReply(&orderapi.BeginReviseOrderReply{RevisedOrderTotal: newTotal}).Success()}, nil
}

func (h CommandHandlers) UndoRevise(ctx context.Context, cmdMsg saga.Command) ([]msg.Reply, error) {
	cmd := cmdMsg.Command().(orderapi.UndoReviseOrder)

	err := h.app.Commands.UndoReviseOrder.Handle(ctx, commands.UndoReviseOrder{OrderID: cmd.OrderID})
	if err != nil {
		return []msg.Reply{msg.WithFailure()}, nil
	}

	return []msg.Reply{msg.WithSuccess()}, nil
}

func (h CommandHandlers) ConfirmRevise(ctx context.Context, cmdMsg saga.Command) ([]msg.Reply, error) {
	cmd := cmdMsg.Command().(orderapi.ConfirmReviseOrder)

	err := h.app.Commands.ConfirmReviseOrder.Handle(ctx, commands.ConfirmReviseOrder{
		OrderID:           cmd.OrderID,
		RevisedQuantities: cmd.RevisedQuantities,
	})
	if err != nil {
		return []msg.Reply{msg.WithFailure()}, nil
	}

	return []msg.Reply{msg.WithSuccess()}, nil
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
		Menu:         evt.Menu,
	})
}

func (h RestaurantEventHandlers) RestaurantMenuRevised(ctx context.Context, evtMsg msg.EntityEvent) error {
	evt := evtMsg.Event().(*restaurantapi.RestaurantMenuRevised)

	return h.app.Commands.ReviseRestaurantMenu.Handle(ctx, commands.ReviseRestaurantMenu{
		RestaurantID: evtMsg.EntityID(),
		Menu:         evt.Menu,
	})
}

type OrderEventHandlers struct{ app Application }

func NewOrderEventHandlers(app Application) OrderEventHandlers { return OrderEventHandlers{app: app} }

func (h OrderEventHandlers) OrderCreated(ctx context.Context, evtMsg msg.EntityEvent) error {
	evt := evtMsg.Event().(*orderapi.OrderCreated)

	err := h.app.Commands.StartCreateOrderSaga.Handle(ctx, commands.StartCreateOrderSaga{
		OrderID:      evtMsg.EntityID(),
		ConsumerID:   evt.ConsumerID,
		RestaurantID: evt.RestaurantID,
		LineItems:    evt.LineItems,
		OrderTotal:   evt.OrderTotal,
	})
	if err != nil {
		return err
	}

	h.app.Metrics.OrdersPlaced.Inc()

	return nil
}

type WebHandlers struct {
	app Application
}

func NewWebHandlers(app Application) WebHandlers {
	return WebHandlers{app: app}
}

func (s WebHandlers) CreateOrder(w http.ResponseWriter, r *http.Request) {
	request := CreateOrderJSONRequestBody{}
	if err := render.Decode(r, &request); err != nil {
		render.Render(w, r, web.NewErrorResponse(err))
		return
	}

	orderID, err := s.app.Commands.CreateOrder.Handle(r.Context(), commands.CreateOrder{
		ConsumerID:   request.ConsumerId,
		RestaurantID: request.RestaurantId,
		DeliverAt:    request.DeliveryTime,
		DeliverTo:    request.DeliveryAddress,
		LineItems:    request.LineItems,
	})
	if err != nil {
		render.Render(w, r, web.NewErrorResponse(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.Respond(w, r, OrderIDResponse{Id: orderID})
}

func (s WebHandlers) GetOrder(w http.ResponseWriter, r *http.Request, orderID OrderID) {
	oid := string(orderID)

	order, err := s.app.Queries.GetOrder.Handle(r.Context(), queries.GetOrder{OrderID: oid})
	if err != nil {
		render.Render(w, r, web.NewErrorResponse(err))
		return
	}

	render.Respond(w, r, OrderResponse{
		OrderId:    oid,
		OrderTotal: order.OrderTotal(),
		State:      order.State.String(),
	})
}

func (s WebHandlers) CancelOrder(w http.ResponseWriter, r *http.Request, orderID OrderID) {
	oid := string(orderID)

	status, err := s.app.Commands.StartCancelOrderSaga.Handle(r.Context(), commands.StartCancelOrderSaga{OrderID: oid})
	if err != nil {
		render.Render(w, r, web.NewErrorResponse(err))
		return
	}

	render.Status(r, http.StatusAccepted)
	render.Respond(w, r, OrderStatusResponse{Status: status})
}

func (s WebHandlers) ReviseOrder(w http.ResponseWriter, r *http.Request, orderID OrderID) {
	oid := string(orderID)

	request := ReviseOrderJSONRequestBody{}
	if err := render.Decode(r, &request); err != nil {
		render.Render(w, r, web.NewErrorResponse(err))
		return
	}

	status, err := s.app.Commands.StartReviseOrderSaga.Handle(r.Context(), commands.StartReviseOrderSaga{
		OrderID:           oid,
		RevisedQuantities: request.RevisedQuantities,
	})
	if err != nil {
		render.Render(w, r, web.NewErrorResponse(err))
		return
	}

	render.Status(r, http.StatusAccepted)
	render.Respond(w, r, OrderStatusResponse{Status: status})
}

func (s WebHandlers) GetRestaurant(w http.ResponseWriter, r *http.Request, restaurantID RestaurantID) {
	rid := string(restaurantID)

	restaurant, err := s.app.Queries.GetRestaurant.Handle(r.Context(), queries.GetRestaurant{RestaurantID: rid})
	if err != nil {
		render.Render(w, r, web.NewErrorResponse(err))
		return
	}

	render.Respond(w, r, RestaurantIDResponse{Id: restaurant.RestaurantID})
}
