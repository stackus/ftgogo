package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/stackus/edat/msg"
	"github.com/stackus/edat/saga"

	"github.com/stackus/ftgogo/kitchen/internal/adapters"
	"github.com/stackus/ftgogo/kitchen/internal/application/commands"
	"github.com/stackus/ftgogo/kitchen/internal/application/queries"
	"github.com/stackus/ftgogo/kitchen/internal/domain"
	"serviceapis"
	"serviceapis/kitchenapi"
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

func initApplication(svc *applications.Service) error {
	serviceapis.RegisterTypes()
	domain.RegisterTypes()

	ticketRepo := adapters.NewTicketRepository(svc.AggregateStore)
	ticketPublisher := adapters.NewTicketPublisher(svc.Publisher)
	restaurantRepo := adapters.NewRestaurantPostgresRepository(svc.PgConn)

	application := Application{
		Commands: Commands{
			CreateTicket:         commands.NewCreateTicketHandler(ticketRepo),
			ConfirmCreateTicket:  commands.NewConfirmCreateTicketHandler(ticketRepo, ticketPublisher),
			CancelCreateTicket:   commands.NewCancelCreateTicketHandler(ticketRepo),
			BeginCancelTicket:    commands.NewBeginCancelTicketHandler(ticketRepo),
			ConfirmCancelTicket:  commands.NewConfirmCancelTicketHandler(ticketRepo, ticketPublisher),
			UndoCancelTicket:     commands.NewUndoCancelTicketHandler(ticketRepo),
			BeginReviseTicket:    commands.NewBeginReviseTicketHandler(ticketRepo),
			ConfirmReviseTicket:  commands.NewConfirmReviseTicketHandler(ticketRepo, ticketPublisher),
			UndoReviseTicket:     commands.NewUndoReviseTicketHandler(ticketRepo),
			AcceptTicket:         commands.NewAcceptTicketHandler(ticketRepo, ticketPublisher),
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

	restaurantEventHandlers := NewRestaurantEventHandlers(application)
	svc.Subscriber.Subscribe(restaurantapi.RestaurantAggregateChannel, msg.NewEntityEventDispatcher().
		Handle(restaurantapi.RestaurantCreated{}, restaurantEventHandlers.RestaurantCreated).
		Handle(restaurantapi.RestaurantMenuRevised{}, restaurantEventHandlers.RestaurantMenuRevised))

	kitchenapi.RegisterKitchenServiceServer(svc.RpcServer, newRpcHandlers(application))

	svc.WebServer.Mount(svc.Cfg.Web.ApiPath, func(r chi.Router) http.Handler {
		return HandlerFromMux(NewWebHandlers(application), r)
	})

	return nil
}

type WebHandlers struct{ app Application }

func NewWebHandlers(app Application) WebHandlers { return WebHandlers{app: app} }

func (h WebHandlers) GetRestaurant(w http.ResponseWriter, r *http.Request, restaurantID RestaurantID) {
	rid := string(restaurantID)

	_, err := h.app.Queries.GetRestaurant.Handle(r.Context(), queries.GetRestaurant{RestaurantID: rid})
	if err != nil {
		render.Render(w, r, web.NewErrorResponse(err))
		return
	}

	render.Respond(w, r, RestaurantIDResponse{Id: rid})
}

func (h WebHandlers) AcceptTicket(w http.ResponseWriter, r *http.Request, ticketID TicketID) {
	tid := string(ticketID)

	request := AcceptTicketJSONRequestBody{}

	if err := render.Decode(r, &request); err != nil {
		render.Render(w, r, web.NewErrorResponse(err))
		return
	}

	err := h.app.Commands.AcceptTicket.Handle(r.Context(), commands.AcceptTicket{
		TicketID: tid,
		ReadyBy:  request.ReadyBy,
	})
	if err != nil {
		render.Render(w, r, web.NewErrorResponse(err))
		return
	}

	render.Status(r, http.StatusAccepted)
	render.Respond(w, r, TicketIDResponse{Id: tid})
}

type RpcHandlers struct {
	app Application
	kitchenapi.UnimplementedKitchenServiceServer
}

var _ kitchenapi.KitchenServiceServer = (*RpcHandlers)(nil)

func newRpcHandlers(app Application) RpcHandlers {
	return RpcHandlers{app: app}
}

func (h RpcHandlers) GetRestaurant(ctx context.Context, request *kitchenapi.GetRestaurantRequest) (*kitchenapi.GetRestaurantResponse, error) {
	_, err := h.app.Queries.GetRestaurant.Handle(ctx, queries.GetRestaurant{RestaurantID: request.RestaurantID})
	if err != nil {
		return nil, err
	}

	return &kitchenapi.GetRestaurantResponse{RestaurantID: request.RestaurantID}, nil
}

func (h RpcHandlers) AcceptTicket(ctx context.Context, request *kitchenapi.AcceptTicketRequest) (*kitchenapi.AcceptTicketResponse, error) {
	err := h.app.Commands.AcceptTicket.Handle(ctx, commands.AcceptTicket{
		TicketID: request.TicketID,
		ReadyBy:  request.ReadyBy.AsTime(),
	})
	if err != nil {
		return nil, err
	}

	return &kitchenapi.AcceptTicketResponse{TicketID: request.TicketID}, nil
}

type CommandHandlers struct{ app Application }

func NewCommandHandlers(app Application) CommandHandlers { return CommandHandlers{app: app} }

func (h CommandHandlers) CreateTicket(ctx context.Context, cmdMsg saga.Command) ([]msg.Reply, error) {
	cmd := cmdMsg.Command().(*kitchenapi.CreateTicket)

	ticketID, err := h.app.Commands.CreateTicket.Handle(ctx, commands.CreateTicket{
		OrderID:      cmd.OrderID,
		RestaurantID: cmd.RestaurantID,
		LineItems:    cmd.TicketDetails,
	})
	if err != nil {
		return []msg.Reply{msg.WithFailure()}, nil
	}

	return []msg.Reply{msg.WithReply(&kitchenapi.CreateTicketReply{TicketID: ticketID}).Success()}, nil
}

func (h CommandHandlers) ConfirmCreateTicket(ctx context.Context, cmdMsg saga.Command) ([]msg.Reply, error) {
	cmd := cmdMsg.Command().(*kitchenapi.ConfirmCreateTicket)

	err := h.app.Commands.ConfirmCreateTicket.Handle(ctx, commands.ConfirmCreateTicket{TicketID: cmd.TicketID})
	if err != nil {
		return []msg.Reply{msg.WithFailure()}, nil
	}

	return []msg.Reply{msg.WithSuccess()}, nil
}

func (h CommandHandlers) CancelCreateTicket(ctx context.Context, cmdMsg saga.Command) ([]msg.Reply, error) {
	cmd := cmdMsg.Command().(*kitchenapi.CancelCreateTicket)

	err := h.app.Commands.CancelCreateTicket.Handle(ctx, commands.CancelCreateTicket{TicketID: cmd.TicketID})
	if err != nil {
		return []msg.Reply{msg.WithFailure()}, nil
	}

	return []msg.Reply{msg.WithSuccess()}, nil
}

func (h CommandHandlers) BeginCancelTicket(ctx context.Context, cmdMsg saga.Command) ([]msg.Reply, error) {
	cmd := cmdMsg.Command().(*kitchenapi.BeginCancelTicket)

	err := h.app.Commands.BeginCancelTicket.Handle(ctx, commands.BeginCancelTicket{
		TicketID:     cmd.TicketID,
		RestaurantID: cmd.RestaurantID,
	})
	if err != nil {
		return []msg.Reply{msg.WithFailure()}, nil
	}

	return []msg.Reply{msg.WithSuccess()}, nil
}

func (h CommandHandlers) ConfirmCancelTicket(ctx context.Context, cmdMsg saga.Command) ([]msg.Reply, error) {
	cmd := cmdMsg.Command().(*kitchenapi.ConfirmCancelTicket)

	err := h.app.Commands.ConfirmCancelTicket.Handle(ctx, commands.ConfirmCancelTicket{
		TicketID:     cmd.TicketID,
		RestaurantID: cmd.RestaurantID,
	})
	if err != nil {
		return []msg.Reply{msg.WithFailure()}, nil
	}

	return []msg.Reply{msg.WithSuccess()}, nil
}

func (h CommandHandlers) UndoCancelTicket(ctx context.Context, cmdMsg saga.Command) ([]msg.Reply, error) {
	cmd := cmdMsg.Command().(*kitchenapi.UndoCancelTicket)

	err := h.app.Commands.UndoCancelTicket.Handle(ctx, commands.UndoCancelTicket{
		TicketID:     cmd.TicketID,
		RestaurantID: cmd.RestaurantID,
	})
	if err != nil {
		return []msg.Reply{msg.WithFailure()}, nil
	}

	return []msg.Reply{msg.WithSuccess()}, nil
}

func (h CommandHandlers) BeginReviseTicket(ctx context.Context, cmdMsg saga.Command) ([]msg.Reply, error) {
	cmd := cmdMsg.Command().(*kitchenapi.BeginReviseTicket)

	err := h.app.Commands.BeginReviseTicket.Handle(ctx, commands.BeginReviseTicket{
		TicketID:          cmd.TicketID,
		RestaurantID:      cmd.RestaurantID,
		RevisedQuantities: cmd.RevisedQuantities,
	})
	if err != nil {
		return []msg.Reply{msg.WithFailure()}, nil
	}

	return []msg.Reply{msg.WithSuccess()}, nil
}

func (h CommandHandlers) ConfirmReviseTicket(ctx context.Context, cmdMsg saga.Command) ([]msg.Reply, error) {
	cmd := cmdMsg.Command().(*kitchenapi.ConfirmReviseTicket)

	err := h.app.Commands.ConfirmReviseTicket.Handle(ctx, commands.ConfirmReviseTicket{
		TicketID:          cmd.TicketID,
		RestaurantID:      cmd.RestaurantID,
		RevisedQuantities: cmd.RevisedQuantities,
	})
	if err != nil {
		return []msg.Reply{msg.WithFailure()}, nil
	}

	return []msg.Reply{msg.WithSuccess()}, nil
}

func (h CommandHandlers) UndoReviseTicket(ctx context.Context, cmdMsg saga.Command) ([]msg.Reply, error) {
	cmd := cmdMsg.Command().(*kitchenapi.UndoReviseTicket)

	err := h.app.Commands.UndoReviseTicket.Handle(ctx, commands.UndoReviseTicket{
		TicketID:     cmd.TicketID,
		RestaurantID: cmd.RestaurantID,
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
