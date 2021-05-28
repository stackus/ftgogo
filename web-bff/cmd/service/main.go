package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"google.golang.org/grpc"

	"github.com/stackus/ftgogo/serviceapis/consumerapi/pb"
	"github.com/stackus/ftgogo/serviceapis/orderapi/pb"
	"github.com/stackus/ftgogo/serviceapis/orderhistoryapi/pb"
	"github.com/stackus/ftgogo/web-bff/internal/adapters"
	"github.com/stackus/ftgogo/web-bff/internal/application"
	"github.com/stackus/ftgogo/web-bff/internal/application/commands"
	"github.com/stackus/ftgogo/web-bff/internal/application/queries"
	"github.com/stackus/ftgogo/web-bff/internal/ports"
	"shared-go/applications"
	"shared-go/rpc"
)

func main() {
	svc := applications.NewBFF(initBFF)
	if err := svc.Execute(); err != nil {
		panic(err)
	}
}

func initBFF(bff *applications.BFF) error {
	var err error

	var orderConn *grpc.ClientConn
	orderConn, err = rpc.NewClientConn(bff.Cfg.Rpc, "orderservice:8000", rpc.WithClientUnaryConvertStatus())
	if err != nil {
		return err
	}
	orderClient := adapters.NewOrderGrpcClient(orderpb.NewOrderServiceClient(orderConn))

	var consumerConn *grpc.ClientConn
	consumerConn, err = rpc.NewClientConn(bff.Cfg.Rpc, "consumerservice:8000", rpc.WithClientUnaryConvertStatus())
	if err != nil {
		return err
	}
	consumerClient := adapters.NewConsumerGrpcClient(consumerpb.NewConsumerServiceClient(consumerConn))

	var orderHistoryConn *grpc.ClientConn
	orderHistoryConn, err = rpc.NewClientConn(bff.Cfg.Rpc, "orderhistoryservice:8000", rpc.WithClientUnaryConvertStatus())
	if err != nil {
		return err
	}
	orderHistoryClient := adapters.NewOrderHistoryGrpcClient(orderhistorypb.NewOrderHistoryServiceClient(orderHistoryConn))

	app := application.Application{
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

	bff.WebServer.Mount(bff.Cfg.Web.ApiPath, ports.NewWebHandlers(app).Mount)

	bff.WebServer.Mount("/spec", func(router chi.Router) http.Handler {
		router.Get("/", func(w http.ResponseWriter, r *http.Request) {
			swagger, err := ports.GetSwagger()
			if err != nil {
				render.Status(r, http.StatusInternalServerError)
				render.PlainText(w, r, fmt.Sprintf("Error rendering Swagger API: %s", err.Error()))
				return
			}
			render.JSON(w, r, swagger)
		})
		return router
	})

	// Perform custom service cleanup
	bff.Cleanup = func(ctx context.Context) error {
		// close all grpc connections; TODO error handling
		_ = orderConn.Close()
		_ = consumerConn.Close()
		_ = orderHistoryConn.Close()

		return nil
	}

	return nil
}
