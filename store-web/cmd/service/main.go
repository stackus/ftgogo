package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	accountingpb "github.com/stackus/ftgogo/serviceapis/accountingapi/pb"
	deliverypb "github.com/stackus/ftgogo/serviceapis/deliveryapi/pb"
	"github.com/stackus/ftgogo/store-web/internal/application/commands"
	"github.com/stackus/ftgogo/store-web/internal/application/queries"
	"shared-go/applications"

	consumerpb "github.com/stackus/ftgogo/serviceapis/consumerapi/pb"
	orderpb "github.com/stackus/ftgogo/serviceapis/orderapi/pb"
	"github.com/stackus/ftgogo/store-web/internal/adapters"
	"github.com/stackus/ftgogo/store-web/internal/application"
	"github.com/stackus/ftgogo/store-web/internal/handlers"
)

func main() {
	gateway := applications.NewGateway(initGateway)
	if err := gateway.Execute(); err != nil {
		panic(err)
	}
}

func initGateway(gateway *applications.Gateway) error {
	// Driven
	accountingClient := adapters.NewAccountingGrpcRepository(accountingpb.NewAccountingServiceClient(gateway.Clients[applications.AccountingService]))
	consumerClient := adapters.NewConsumerGrpcClient(consumerpb.NewConsumerServiceClient(gateway.Clients[applications.ConsumerService]))
	deliveryClient := adapters.NewDeliveryGrpcRepository(deliverypb.NewDeliveryServiceClient(gateway.Clients[applications.DeliveryService]))
	orderClient := adapters.NewOrderGrpcClient(orderpb.NewOrderServiceClient(gateway.Clients[applications.OrderService]))

	app := application.Service{
		Commands: application.Commands{
			EnableAccount:          commands.NewEnableAccountHandler(accountingClient),
			DisableAccount:         commands.NewDisableAccountHandler(accountingClient),
			SetCourierAvailability: commands.NewSetCourierAvailabilityHandler(deliveryClient),
			CancelOrder:            commands.NewCancelOrderHandler(orderClient),
		},
		Queries: application.Queries{
			GetAccount:         queries.NewGetAccountHandler(accountingClient),
			GetConsumer:        queries.NewGetConsumerHandler(consumerClient),
			GetDeliveryHistory: queries.NewGetDeliveryHistoryHandler(deliveryClient),
			GetOrder:           queries.NewGetOrderHandler(orderClient),
		},
	}

	// Drivers
	gateway.WebServer.Mount(gateway.Cfg.Web.ApiPath, handlers.NewWebHandlers(app).Mount)

	gateway.WebServer.Mount("/spec", func(router chi.Router) http.Handler {
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
