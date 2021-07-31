package handlers

import (
	"context"

	"google.golang.org/grpc"

	"github.com/stackus/ftgogo/kitchen/internal/application"
	"github.com/stackus/ftgogo/kitchen/internal/application/commands"
	"github.com/stackus/ftgogo/kitchen/internal/application/queries"
	"github.com/stackus/ftgogo/serviceapis/kitchenapi/pb"
)

type RpcHandlers struct {
	app application.ServiceApplication
	kitchenpb.UnimplementedKitchenServiceServer
}

var _ kitchenpb.KitchenServiceServer = (*RpcHandlers)(nil)

func NewRpcHandlers(app application.ServiceApplication) RpcHandlers {
	return RpcHandlers{app: app}
}

func (h RpcHandlers) Mount(registrar grpc.ServiceRegistrar) {
	kitchenpb.RegisterKitchenServiceServer(registrar, h)
}

func (h RpcHandlers) GetRestaurant(ctx context.Context, request *kitchenpb.GetRestaurantRequest) (*kitchenpb.GetRestaurantResponse, error) {
	_, err := h.app.GetRestaurant(ctx, queries.GetRestaurant{RestaurantID: request.RestaurantID})
	if err != nil {
		return nil, err
	}

	return &kitchenpb.GetRestaurantResponse{RestaurantID: request.RestaurantID}, nil
}

func (h RpcHandlers) AcceptTicket(ctx context.Context, request *kitchenpb.AcceptTicketRequest) (*kitchenpb.AcceptTicketResponse, error) {
	err := h.app.AcceptTicket(ctx, commands.AcceptTicket{
		TicketID: request.TicketID,
		ReadyBy:  request.ReadyBy.AsTime(),
	})
	if err != nil {
		return nil, err
	}

	return &kitchenpb.AcceptTicketResponse{TicketID: request.TicketID}, nil
}
