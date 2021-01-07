package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"github.com/stackus/edat/msg"
	"github.com/stackus/edat/saga"
	"github.com/stackus/ftgogo/consumer/internal/adapters"
	"github.com/stackus/ftgogo/consumer/internal/application/commands"
	"github.com/stackus/ftgogo/consumer/internal/application/queries"
	"github.com/stackus/ftgogo/consumer/internal/domain"
	"serviceapis"
	"serviceapis/consumerapi"
	"shared-go/applications"
	"shared-go/errs"
	"shared-go/web"
)

// To regenerate the web server api use the following generate command
//go:generate oapi-codegen -config oapi-codegen.cfg.yaml openapi.yaml

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	RegisterConsumer        commands.RegisterConsumerHandler
	ValidateOrderByConsumer commands.ValidateOrderByConsumerHandler
}

type Queries struct {
	GetConsumer queries.GetConsumerHandler
}

func initApplication(svc *applications.Service) error {
	serviceapis.RegisterTypes()
	domain.RegisterTypes()

	consumerRepo := adapters.NewConsumerRepository(svc.AggregateStore)
	consumerPublisher := adapters.NewConsumerPublisher(svc.Publisher)

	application := Application{
		Commands: Commands{
			RegisterConsumer:        commands.NewRegisterConsumerHandler(consumerRepo, consumerPublisher),
			ValidateOrderByConsumer: commands.NewValidateOrderByConsumerHandler(consumerRepo),
		},
		Queries: Queries{
			GetConsumer: queries.NewGetConsumerHandler(consumerRepo),
		},
	}

	cmdHandlers := NewCommandHandlers(application)
	svc.Subscriber.Subscribe(consumerapi.ConsumerServiceCommandChannel, saga.NewCommandDispatcher(svc.Publisher).
		Handle(consumerapi.ValidateOrderByConsumer{}, cmdHandlers.ValidateOrderByConsumer))

	// TODO refactor so a string isn't used here
	svc.WebServer.Mount("/api", func(r chi.Router) http.Handler {
		return HandlerFromMux(NewWebHandlers(application), r)
	})

	return nil
}

type WebHandlers struct{ app Application }

func NewWebHandlers(app Application) WebHandlers { return WebHandlers{app: app} }

func (h WebHandlers) RegisterConsumer(w http.ResponseWriter, r *http.Request) {
	request := RegisterConsumerJSONRequestBody{}

	if err := render.Decode(r, &request); err != nil {
		render.Render(w, r, web.NewErrorResponse(err))
		return
	}

	consumerID, err := h.app.Commands.RegisterConsumer.Handle(r.Context(), commands.RegisterConsumer{
		Name: request.Name,
	})
	if err != nil {
		render.Render(w, r, web.NewErrorResponse(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.Respond(w, r, ConsumerIDResponse{Id: consumerID})
}

func (h WebHandlers) GetConsumer(w http.ResponseWriter, r *http.Request, consumerID ConsumerID) {
	cid := string(consumerID)

	consumer, err := h.app.Queries.GetConsumer.Handle(r.Context(), queries.GetConsumer{ConsumerID: cid})
	if err != nil {
		render.Render(w, r, web.NewErrorResponse(err))
		return
	}

	render.Respond(w, r, ConsumerResponse{
		ConsumerId: consumer.ID(),
		Name:       consumer.Name(),
	})
}

type CommandHandlers struct{ app Application }

func NewCommandHandlers(app Application) CommandHandlers { return CommandHandlers{app: app} }

func (h CommandHandlers) ValidateOrderByConsumer(ctx context.Context, cmdMsg saga.Command) ([]msg.Reply, error) {
	cmd := cmdMsg.Command().(*consumerapi.ValidateOrderByConsumer)

	err := h.app.Commands.ValidateOrderByConsumer.Handle(ctx, commands.ValidateOrderByConsumer{
		ConsumerID: cmd.ConsumerID,
		OrderID:    cmd.OrderID,
		OrderTotal: cmd.OrderTotal,
	})
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return []msg.Reply{msg.WithReply(&consumerapi.ConsumerNotFound{}).Failure()}, nil
		}

		return []msg.Reply{msg.WithFailure()}, nil
	}

	return []msg.Reply{msg.WithSuccess()}, nil
}
