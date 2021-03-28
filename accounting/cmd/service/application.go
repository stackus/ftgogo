package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/stackus/edat/msg"
	"github.com/stackus/edat/saga"

	"github.com/stackus/ftgogo/account/internal/adapters"
	"github.com/stackus/ftgogo/account/internal/application/commands"
	"github.com/stackus/ftgogo/account/internal/application/queries"
	"github.com/stackus/ftgogo/account/internal/domain"
	"serviceapis"
	"serviceapis/accountingapi"
	"serviceapis/consumerapi"
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
	AuthorizeOrder        commands.AuthorizeOrderHandler
	ReverseAuthorizeOrder commands.ReverseAuthorizeOrderHandler
	ReviseAuthorizeOrder  commands.ReviseAuthorizeOrderHandler
	CreateAccount         commands.CreateAccountHandler
	DisableAccount        commands.DisableAccountHandler
	EnableAccount         commands.EnableAccountHandler
}

type Queries struct {
	GetAccount queries.GetAccountHandler
}

func initApplication(svc *applications.Service) error {
	serviceapis.RegisterTypes()
	domain.RegisterTypes()

	accountRepo := adapters.NewAccountRepository(svc.AggregateStore)

	application := Application{
		Commands: Commands{
			AuthorizeOrder:        commands.NewAuthorizeOrderHandler(accountRepo),
			ReverseAuthorizeOrder: commands.NewReverseAuthorizeOrderHandler(accountRepo),
			ReviseAuthorizeOrder:  commands.NewReviseAuthorizeOrderHandler(accountRepo),
			CreateAccount:         commands.NewCreateAccountHandler(accountRepo),
			DisableAccount:        commands.NewDisableAccountHandler(accountRepo),
			EnableAccount:         commands.NewEnableAccountHandler(accountRepo),
		},
		Queries: Queries{
			GetAccount: queries.NewGetAccountHandler(accountRepo),
		},
	}

	cmdHandlers := NewCommandHandlers(application)
	svc.Subscriber.Subscribe(accountingapi.AccountingServiceCommandChannel, saga.NewCommandDispatcher(svc.Publisher).
		Handle(accountingapi.AuthorizeOrder{}, cmdHandlers.AuthorizeOrder).
		Handle(accountingapi.ReverseAuthorizeOrder{}, cmdHandlers.ReverseAuthorizeOrder).
		Handle(accountingapi.ReviseAuthorizeOrder{}, cmdHandlers.ReviseAuthorizeOrder))

	consumerEventHandlers := NewConsumerEventHandlers(application)
	svc.Subscriber.Subscribe(consumerapi.ConsumerAggregateChannel, msg.NewEntityEventDispatcher().
		Handle(consumerapi.ConsumerRegistered{}, consumerEventHandlers.ConsumerRegistered))

	accountingapi.RegisterAccountingServiceServer(svc.RpcServer, newRpcHandlers(application))

	svc.WebServer.Mount(svc.Cfg.Web.ApiPath, func(r chi.Router) http.Handler {
		return HandlerFromMux(NewWebHandlers(application), r)
	})

	return nil
}

type WebHandlers struct {
	app Application
}

func NewWebHandlers(app Application) WebHandlers {
	return WebHandlers{app: app}
}

func (h WebHandlers) GetAccount(w http.ResponseWriter, r *http.Request, accountID AccountID) {
	aid := string(accountID)

	_, err := h.app.Queries.GetAccount.Handle(r.Context(), queries.GetAccount{AccountID: aid})
	if err != nil {
		render.Render(w, r, web.NewErrorResponse(err))
		return
	}

	render.Respond(w, r, AccountIDResponse{Id: aid})
}

func (h WebHandlers) DisableAccount(w http.ResponseWriter, r *http.Request, accountID AccountID) {
	aid := string(accountID)

	err := h.app.Commands.DisableAccount.Handle(r.Context(), commands.DisableAccount{AccountID: aid})
	if err != nil {
		render.Render(w, r, web.NewErrorResponse(err))
		return
	}

	render.Respond(w, r, AccountIDResponse{Id: aid})
}

func (h WebHandlers) EnableAccount(w http.ResponseWriter, r *http.Request, accountID AccountID) {
	aid := string(accountID)

	err := h.app.Commands.EnableAccount.Handle(r.Context(), commands.EnableAccount{AccountID: aid})
	if err != nil {
		render.Render(w, r, web.NewErrorResponse(err))
		return
	}

	render.Respond(w, r, AccountIDResponse{Id: aid})
}

type RpcHandlers struct {
	app Application
	accountingapi.UnimplementedAccountingServiceServer
}

var _ accountingapi.AccountingServiceServer = (*RpcHandlers)(nil)

func newRpcHandlers(app Application) RpcHandlers {
	return RpcHandlers{app: app}
}

func (h RpcHandlers) GetAccount(ctx context.Context, request *accountingapi.GetAccountRequest) (*accountingapi.GetAccountResponse, error) {
	_, err := h.app.Queries.GetAccount.Handle(ctx, queries.GetAccount{AccountID: request.AccountID})
	if err != nil {
		return nil, err
	}

	return &accountingapi.GetAccountResponse{AccountID: request.AccountID}, nil
}

func (h RpcHandlers) DisableAccount(ctx context.Context, request *accountingapi.DisableAccountRequest) (*accountingapi.DisableAccountResponse, error) {
	err := h.app.Commands.DisableAccount.Handle(ctx, commands.DisableAccount{AccountID: request.AccountID})
	if err != nil {
		return nil, err
	}

	return &accountingapi.DisableAccountResponse{AccountID: request.AccountID}, nil
}

func (h RpcHandlers) EnableAccount(ctx context.Context, request *accountingapi.EnableAccountRequest) (*accountingapi.EnableAccountResponse, error) {
	err := h.app.Commands.EnableAccount.Handle(ctx, commands.EnableAccount{AccountID: request.AccountID})
	if err != nil {
		return nil, err
	}

	return &accountingapi.EnableAccountResponse{AccountID: request.AccountID}, nil
}

type CommandHandlers struct {
	app Application
}

func NewCommandHandlers(app Application) CommandHandlers {
	return CommandHandlers{app: app}
}

func (h CommandHandlers) AuthorizeOrder(ctx context.Context, cmdMsg saga.Command) ([]msg.Reply, error) {
	cmd := cmdMsg.Command().(*accountingapi.AuthorizeOrder)

	err := h.app.Commands.AuthorizeOrder.Handle(ctx, commands.AuthorizeOrder{
		ConsumerID: cmd.ConsumerID,
		OrderID:    cmd.OrderID,
		OrderTotal: cmd.OrderTotal,
	})
	if err != nil {
		if errors.Is(err, domain.ErrAccountDisabled) {
			return []msg.Reply{msg.WithReply(&accountingapi.AccountDisabled{}).Failure()}, nil
		}
		return nil, err
	}

	return []msg.Reply{msg.WithSuccess()}, nil
}

func (h CommandHandlers) ReverseAuthorizeOrder(ctx context.Context, cmdMsg saga.Command) ([]msg.Reply, error) {
	cmd := cmdMsg.Command().(*accountingapi.ReverseAuthorizeOrder)

	err := h.app.Commands.ReverseAuthorizeOrder.Handle(ctx, commands.ReverseAuthorizeOrder{
		ConsumerID: cmd.ConsumerID,
		OrderID:    cmd.OrderID,
		OrderTotal: cmd.OrderTotal,
	})
	if err != nil {
		if errors.Is(err, domain.ErrAccountDisabled) {
			return []msg.Reply{msg.WithReply(&accountingapi.AccountDisabled{}).Failure()}, nil
		}
		return nil, err
	}

	return []msg.Reply{msg.WithSuccess()}, nil
}

func (h CommandHandlers) ReviseAuthorizeOrder(ctx context.Context, cmdMsg saga.Command) ([]msg.Reply, error) {
	cmd := cmdMsg.Command().(*accountingapi.ReviseAuthorizeOrder)

	err := h.app.Commands.ReviseAuthorizeOrder.Handle(ctx, commands.ReviseAuthorizeOrder{
		ConsumerID: cmd.ConsumerID,
		OrderID:    cmd.OrderID,
		OrderTotal: cmd.OrderTotal,
	})
	if err != nil {
		if errors.Is(err, domain.ErrAccountDisabled) {
			return []msg.Reply{msg.WithReply(&accountingapi.AccountDisabled{}).Failure()}, nil
		}
		return nil, err
	}

	return []msg.Reply{msg.WithSuccess()}, nil
}

type ConsumerEventHandlers struct {
	app Application
}

func NewConsumerEventHandlers(app Application) ConsumerEventHandlers {
	return ConsumerEventHandlers{app: app}
}

func (h ConsumerEventHandlers) ConsumerRegistered(ctx context.Context, evtMsg msg.EntityEvent) error {
	evt := evtMsg.Event().(*consumerapi.ConsumerRegistered)

	return h.app.Commands.CreateAccount.Handle(ctx, commands.CreateAccount{
		ConsumerID: evtMsg.EntityID(),
		Name:       evt.Name,
	})
}
