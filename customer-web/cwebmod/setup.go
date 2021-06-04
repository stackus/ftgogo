package cwebmod

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/stackus/errors"
	"google.golang.org/grpc"

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
	// Infrastructure
	conn, err := grpc.Dial(svc.Cfg.Rpc.Address, grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
		addr, err := net.ResolveUnixAddr(svc.Cfg.Rpc.Network, svc.Cfg.Rpc.Address)
		if err != nil {
			return nil, err
		}
		return net.DialUnix(svc.Cfg.Rpc.Network, nil, addr)
	}), grpc.WithChainUnaryInterceptor(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return errors.ReceiveGRPCError(invoker(ctx, method, req, reply, cc, opts...))
	}))
	if err != nil {
		return err
	}

	// Driven
	// Create a unix socket connection to the GRPC server
	orderClient := adapters.NewOrderGrpcClient(orderpb.NewOrderServiceClient(conn))
	consumerClient := adapters.NewConsumerGrpcClient(consumerpb.NewConsumerServiceClient(conn))
	orderHistoryClient := adapters.NewOrderHistoryGrpcClient(orderhistorypb.NewOrderHistoryServiceClient(conn))

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

	svc.Cleanup = func(ctx context.Context) error {
		return conn.Close()
	}

	return nil
}
