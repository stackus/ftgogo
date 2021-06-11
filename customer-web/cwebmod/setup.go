package cwebmod

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/stackus/ftgogo/customer-web/internal/adapters"
	"github.com/stackus/ftgogo/customer-web/internal/application"
	"github.com/stackus/ftgogo/customer-web/internal/application/commands"
	"github.com/stackus/ftgogo/customer-web/internal/application/queries"
	"github.com/stackus/ftgogo/customer-web/internal/handlers"
	consumerpb "github.com/stackus/ftgogo/serviceapis/consumerapi/pb"
	orderpb "github.com/stackus/ftgogo/serviceapis/orderapi/pb"
	orderhistorypb "github.com/stackus/ftgogo/serviceapis/orderhistoryapi/pb"
	"shared-go/applications"
)

func Setup(svc *applications.Monolith) error {
	// Driven
	// Create a unix socket connection to the GRPC server
	orderClient := adapters.NewOrderGrpcClient(orderpb.NewOrderServiceClient(svc.Clients[applications.OrderService]))
	consumerClient := adapters.NewConsumerGrpcClient(consumerpb.NewConsumerServiceClient(svc.Clients[applications.ConsumerService]))
	orderHistoryClient := adapters.NewOrderHistoryGrpcClient(orderhistorypb.NewOrderHistoryServiceClient(svc.Clients[applications.OrderHistoryService]))

	app := application.Service{
		Commands: application.Commands{
			RegisterConsumer:      commands.NewRegisterConsumerHandler(consumerClient),
			CreateOrder:           commands.NewCreateOrderHandler(orderClient, consumerClient),
			ReviseOrder:           commands.NewReviseOrderHandler(orderClient),
			CancelOrder:           commands.NewCancelOrderHandler(orderClient),
			AddConsumerAddress:    commands.NewAddConsumerAddressHandler(consumerClient),
			UpdateConsumerAddress: commands.NewUpdateConsumerAddressHandler(consumerClient),
			RemoveConsumerAddress: commands.NewRemoveConsumerAddressHandler(consumerClient),
		},
		Queries: application.Queries{
			GetConsumer:        queries.NewGetConsumerHandler(consumerClient),
			GetOrder:           queries.NewGetOrderHandler(orderClient),
			GetConsumerAddress: queries.NewGetConsumerAddressHandler(consumerClient),
			SearchOrders:       queries.NewSearchOrdersHandler(orderHistoryClient),
		},
	}

	// Drivers
	svc.WebServer.Mount(svc.Cfg.Web.ApiPath, handlers.NewWebHandlers(app).Mount)

	svc.WebServer.Mount("/spec", func(router chi.Router) http.Handler {
		router.Get("/", func(w http.ResponseWriter, r *http.Request) {
			swagger, err := handlers.GetSwagger()
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
