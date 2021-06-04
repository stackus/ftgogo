package handlers

import (
	"context"

	"google.golang.org/grpc"

	"github.com/stackus/ftgogo/restaurant/internal/application"
	"github.com/stackus/ftgogo/restaurant/internal/application/commands"
	"github.com/stackus/ftgogo/restaurant/internal/application/queries"
	"github.com/stackus/ftgogo/serviceapis/commonapi/pb"
	"github.com/stackus/ftgogo/serviceapis/restaurantapi"
	"github.com/stackus/ftgogo/serviceapis/restaurantapi/pb"
)

type RpcHandlers struct {
	app application.Application
	restaurantpb.UnimplementedRestaurantServiceServer
}

var _ restaurantpb.RestaurantServiceServer = (*RpcHandlers)(nil)

func NewRpcHandlers(app application.Application) RpcHandlers {
	return RpcHandlers{app: app}
}

func (h RpcHandlers) Mount(registrar grpc.ServiceRegistrar) {
	restaurantpb.RegisterRestaurantServiceServer(registrar, h)
}

func (h RpcHandlers) CreateRestaurant(ctx context.Context, request *restaurantpb.CreateRestaurantRequest) (*restaurantpb.CreateRestaurantResponse, error) {
	menuItems := make([]restaurantapi.MenuItem, 0)
	for _, item := range request.Menu.MenuItems {
		menuItems = append(menuItems, restaurantapi.MenuItem{
			ID:    item.ID,
			Name:  item.Name,
			Price: int(item.Price),
		})
	}

	restaurantID, err := h.app.Commands.CreateRestaurant.Handle(ctx, commands.CreateRestaurant{
		Name: request.Name,
		Address: restaurantapi.Address{
			Street1: request.Address.Street1,
			Street2: request.Address.Street2,
			City:    request.Address.City,
			State:   request.Address.State,
			Zip:     request.Address.Zip,
		},
		MenuItems: menuItems,
	})
	if err != nil {
		return nil, err
	}

	return &restaurantpb.CreateRestaurantResponse{RestaurantID: restaurantID}, nil
}

func (h RpcHandlers) GetRestaurant(ctx context.Context, request *restaurantpb.GetRestaurantRequest) (*restaurantpb.GetRestaurantResponse, error) {
	restaurant, err := h.app.Queries.GetRestaurant.Handle(ctx, queries.GetRestaurant{RestaurantID: request.RestaurantID})
	if err != nil {
		return nil, err
	}

	menuItems := make([]*restaurantpb.GetRestaurantResponseMenuItem, 0, len(restaurant.MenuItems))
	for _, item := range restaurant.MenuItems {
		menuItems = append(menuItems, &restaurantpb.GetRestaurantResponseMenuItem{
			ID:    item.ID,
			Name:  item.Name,
			Price: int64(item.Price),
		})
	}

	return &restaurantpb.GetRestaurantResponse{
		RestaurantID: restaurant.ID(),
		Name:         restaurant.Name,
		Address: &commonpb.Address{
			Street1: restaurant.Address.Street1,
			Street2: restaurant.Address.Street2,
			City:    restaurant.Address.City,
			State:   restaurant.Address.State,
			Zip:     restaurant.Address.Zip,
		},
		Menu: &restaurantpb.GetRestaurantResponseMenu{MenuItems: menuItems},
	}, nil
}
