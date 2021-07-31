package adapters

import (
	"context"

	"github.com/stackus/ftgogo/serviceapis/commonapi"
	"github.com/stackus/ftgogo/serviceapis/restaurantapi"
	restaurantpb "github.com/stackus/ftgogo/serviceapis/restaurantapi/pb"
	"github.com/stackus/ftgogo/store-web/internal/domain"
)

type RestaurantGrpcRepository struct {
	client restaurantpb.RestaurantServiceClient
}

var _ RestaurantRepository = (*RestaurantGrpcRepository)(nil)

func NewRestaurantGrpcRepository(client restaurantpb.RestaurantServiceClient) *RestaurantGrpcRepository {
	return &RestaurantGrpcRepository{client: client}
}

func (r RestaurantGrpcRepository) Create(ctx context.Context, createRestaurant CreateRestaurant) (string, error) {
	menuItems := make([]*restaurantpb.MenuItem, 0, len(createRestaurant.MenuItems))
	for _, item := range createRestaurant.MenuItems {
		menuItems = append(menuItems, restaurantapi.ToMenuItemProto(item))
	}
	resp, err := r.client.CreateRestaurant(ctx, &restaurantpb.CreateRestaurantRequest{
		Name:    createRestaurant.Name,
		Address: commonapi.ToAddressProto(createRestaurant.Address),
		Menu:    &restaurantpb.Menu{MenuItems: menuItems},
	})
	if err != nil {
		return "", err
	}
	return resp.RestaurantID, nil
}

func (r RestaurantGrpcRepository) Find(ctx context.Context, findRestaurant FindRestaurant) (*domain.Restaurant, error) {
	resp, err := r.client.GetRestaurant(ctx, &restaurantpb.GetRestaurantRequest{RestaurantID: findRestaurant.RestaurantID})
	if err != nil {
		return nil, err
	}

	menuItems := make([]restaurantapi.MenuItem, 0, len(resp.Menu.MenuItems))
	for _, item := range resp.Menu.MenuItems {
		menuItems = append(menuItems, restaurantapi.FromMenuItemProto(item))
	}

	return &domain.Restaurant{
		RestaurantID: resp.RestaurantID,
		Name:         resp.Name,
		Address:      commonapi.FromAddressProto(resp.Address),
		MenuItems:    menuItems,
	}, nil
}
