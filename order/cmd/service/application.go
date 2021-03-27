package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
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

	ordersPlacedCounter := adapters.NewOrdersPlacedCounter()
	ordersApprovedCounter := adapters.NewOrdersApprovedCounter()
	ordersRejectedCounter := adapters.NewOrdersRejectedCounter()

	application := Application{
		Commands: Commands{
			CreateOrder:          commands.NewCreateOrderHandler(orderRepo, restaurantRepo, orderPublisher, svc.Logger),
			ApproveOrder:         commands.NewApproveOrderHandler(orderRepo, orderPublisher, ordersApprovedCounter),
			RejectOrder:          commands.NewRejectOrderHandler(orderRepo, orderPublisher, ordersRejectedCounter),
			BeginCancelOrder:     commands.NewBeginCancelOrderHandler(orderRepo),
			UndoCancelOrder:      commands.NewUndoCancelOrderHandler(orderRepo),
			ConfirmCancelOrder:   commands.NewConfirmCancelOrderHandler(orderRepo, orderPublisher),
			BeginReviseOrder:     commands.NewBeginReviseOrderHandler(orderRepo, orderPublisher),
			UndoReviseOrder:      commands.NewUndoReviseOrderHandler(orderRepo),
			ConfirmReviseOrder:   commands.NewConfirmReviseOrderHandler(orderRepo, orderPublisher),
			StartCreateOrderSaga: commands.NewStartCreateOrderSagaHandler(createOrderSaga, ordersPlacedCounter),
			StartCancelOrderSaga: commands.NewStartCancelOrderSagaHandler(orderRepo, cancelOrderSaga),
			StartReviseOrderSaga: commands.NewStartReviseOrderSagaHandler(orderRepo, reviseOrderSaga),
			CreateRestaurant:     commands.NewCreateRestaurantHandler(restaurantRepo),
			ReviseRestaurantMenu: commands.NewReviseRestaurantMenuHandler(restaurantRepo),
		},
		Queries: Queries{
			GetOrder:      queries.NewGetOrderHandler(orderRepo),
			GetRestaurant: queries.NewGetRestaurantHandler(restaurantRepo),
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

	orderapi.RegisterOrderServiceServer(svc.RpcServer, newRpcHandlers(application))

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
		DeliverTo: orderapi.Address{
			Street1: request.DeliveryAddress.Street1,
			Street2: request.DeliveryAddress.Street2,
			City:    request.DeliveryAddress.City,
			State:   request.DeliveryAddress.State,
			Zip:     request.DeliveryAddress.Zip,
		},
		LineItems: request.LineItems,
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

type RpcHandlers struct {
	app Application
	orderapi.UnimplementedOrderServiceServer
}

var _ orderapi.OrderServiceServer = (*RpcHandlers)(nil)

func newRpcHandlers(app Application) RpcHandlers {
	return RpcHandlers{app: app}
}

func (h RpcHandlers) CreateOrder(ctx context.Context, request *orderapi.CreateOrderRequest) (*orderapi.CreateOrderResponse, error) {
	lineItems := make(map[string]int, len(request.LineItems))
	for s, i := range request.LineItems {
		lineItems[s] = int(i)
	}

	orderID, err := h.app.Commands.CreateOrder.Handle(ctx, commands.CreateOrder{
		ConsumerID:   request.ConsumerID,
		RestaurantID: request.RestaurantID,
		DeliverAt:    request.DeliverAt.AsTime(),
		DeliverTo: orderapi.Address{
			Street1: request.DeliverTo.Street1,
			Street2: request.DeliverTo.Street2,
			City:    request.DeliverTo.City,
			State:   request.DeliverTo.State,
			Zip:     request.DeliverTo.Zip,
		},
		LineItems: lineItems,
	})
	if err != nil {
		return nil, err
	}

	return &orderapi.CreateOrderResponse{OrderID: orderID}, nil
}

func (h RpcHandlers) GetOrder(ctx context.Context, request *orderapi.GetOrderRequest) (*orderapi.GetOrderResponse, error) {
	order, err := h.app.Queries.GetOrder.Handle(ctx, queries.GetOrder{OrderID: request.OrderID})
	if err != nil {
		return nil, err
	}

	return &orderapi.GetOrderResponse{
		OrderID:    order.ID(),
		OrderTotal: int64(order.OrderTotal()),
		State:      order.State.String(),
	}, nil
}

func (h RpcHandlers) CancelOrder(ctx context.Context, request *orderapi.CancelOrderRequest) (*orderapi.CancelOrderResponse, error) {
	status, err := h.app.Commands.StartCancelOrderSaga.Handle(ctx, commands.StartCancelOrderSaga{OrderID: request.OrderID})
	if err != nil {
		return nil, err
	}

	return &orderapi.CancelOrderResponse{Status: status}, nil
}

func (h RpcHandlers) ReviseOrder(ctx context.Context, request *orderapi.ReviseOrderRequest) (*orderapi.ReviseOrderResponse, error) {
	revisedQuantities := make(map[string]int, len(request.RevisedQuantities))
	for s, i := range request.RevisedQuantities {
		revisedQuantities[s] = int(i)
	}

	status, err := h.app.Commands.StartReviseOrderSaga.Handle(ctx, commands.StartReviseOrderSaga{
		OrderID:           request.OrderID,
		RevisedQuantities: revisedQuantities,
	})
	if err != nil {
		return nil, err
	}

	return &orderapi.ReviseOrderResponse{Status: status}, nil
}

func (h RpcHandlers) GetRestaurant(ctx context.Context, request *orderapi.GetRestaurantRequest) (*orderapi.GetRestaurantResponse, error) {
	restaurant, err := h.app.Queries.GetRestaurant.Handle(ctx, queries.GetRestaurant{RestaurantID: request.RestaurantID})
	if err != nil {
		return nil, err
	}

	return &orderapi.GetRestaurantResponse{RestaurantID: restaurant.RestaurantID}, nil
}
