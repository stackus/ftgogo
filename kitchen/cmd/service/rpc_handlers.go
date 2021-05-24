package main

import (
	"context"

	"github.com/stackus/ftgogo/kitchen/internal/application/commands"
	"github.com/stackus/ftgogo/kitchen/internal/application/queries"
	"github.com/stackus/ftgogo/serviceapis/kitchenapi/pb"
)

type rpcHandlers struct {
	app Application
	kitchenpb.UnimplementedKitchenServiceServer
}

var _ kitchenpb.KitchenServiceServer = (*rpcHandlers)(nil)

func newRpcHandlers(app Application) rpcHandlers {
	return rpcHandlers{app: app}
}

func (h rpcHandlers) GetRestaurant(ctx context.Context, request *kitchenpb.GetRestaurantRequest) (*kitchenpb.GetRestaurantResponse, error) {
	_, err := h.app.Queries.GetRestaurant.Handle(ctx, queries.GetRestaurant{RestaurantID: request.RestaurantID})
	if err != nil {
		return nil, err
	}

	return &kitchenpb.GetRestaurantResponse{RestaurantID: request.RestaurantID}, nil
}

func (h rpcHandlers) AcceptTicket(ctx context.Context, request *kitchenpb.AcceptTicketRequest) (*kitchenpb.AcceptTicketResponse, error) {
	err := h.app.Commands.AcceptTicket.Handle(ctx, commands.AcceptTicket{
		TicketID: request.TicketID,
		ReadyBy:  request.ReadyBy.AsTime(),
	})
	if err != nil {
		return nil, err
	}

	return &kitchenpb.AcceptTicketResponse{TicketID: request.TicketID}, nil
}
