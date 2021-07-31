package handlers

import (
	"context"

	"google.golang.org/grpc"

	"github.com/stackus/ftgogo/restaurant/internal/application"
	"github.com/stackus/ftgogo/restaurant/internal/application/commands"
	"github.com/stackus/ftgogo/restaurant/internal/application/queries"
	"github.com/stackus/ftgogo/serviceapis/commonapi"
	"github.com/stackus/ftgogo/serviceapis/restaurantapi"
	"github.com/stackus/ftgogo/serviceapis/restaurantapi/pb"
)

type RpcHandlers struct {
	app application.ServiceApplication
	restaurantpb.UnimplementedRestaurantServiceServer
}

var _ restaurantpb.RestaurantServiceServer = (*RpcHandlers)(nil)

func NewRpcHandlers(app application.ServiceApplication) RpcHandlers {
	return RpcHandlers{app: app}
}

func (h RpcHandlers) Mount(registrar grpc.ServiceRegistrar) {
	restaurantpb.RegisterRestaurantServiceServer(registrar, h)
}

func (h RpcHandlers) CreateRestaurant(ctx context.Context, request *restaurantpb.CreateRestaurantRequest) (*restaurantpb.CreateRestaurantResponse, error) {
	menuItems := make([]restaurantapi.MenuItem, 0)
	for _, item := range request.Menu.MenuItems {
		menuItems = append(menuItems, restaurantapi.FromMenuItemProto(item))
	}

	restaurantID, err := h.app.CreateRestaurant(ctx, commands.CreateRestaurant{
		Name:      request.Name,
		Address:   commonapi.FromAddressProto(request.Address),
		MenuItems: menuItems,
	})
	if err != nil {
		return nil, err
	}

	return &restaurantpb.CreateRestaurantResponse{RestaurantID: restaurantID}, nil
}

func (h RpcHandlers) GetRestaurant(ctx context.Context, request *restaurantpb.GetRestaurantRequest) (*restaurantpb.GetRestaurantResponse, error) {
	restaurant, err := h.app.GetRestaurant(ctx, queries.GetRestaurant{RestaurantID: request.RestaurantID})
	if err != nil {
		return nil, err
	}

	menuItems := make([]*restaurantpb.MenuItem, 0, len(restaurant.MenuItems))
	for _, item := range restaurant.MenuItems {
		menuItems = append(menuItems, restaurantapi.ToMenuItemProto(item))
	}

	return &restaurantpb.GetRestaurantResponse{
		RestaurantID: restaurant.ID(),
		Name:         restaurant.Name,
		Address:      commonapi.ToAddressProto(restaurant.Address),
		Menu:         &restaurantpb.Menu{MenuItems: menuItems},
	}, nil
}
