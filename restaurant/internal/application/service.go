package application

import (
	"context"

	"github.com/stackus/ftgogo/restaurant/internal/application/commands"
	"github.com/stackus/ftgogo/restaurant/internal/application/ports"
	"github.com/stackus/ftgogo/restaurant/internal/application/queries"
	"github.com/stackus/ftgogo/restaurant/internal/domain"
)

type ServiceApplication interface {
	CreateRestaurant(context.Context, commands.CreateRestaurant) (string, error)
	GetRestaurant(context.Context, queries.GetRestaurant) (*domain.Restaurant, error)
}

type Service struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateRestaurant commands.CreateRestaurantHandler
}

type Queries struct {
	GetRestaurant queries.GetRestaurantHandler
}

func NewServiceApplication(restaurantRepo ports.RestaurantRepository) *Service {
	return &Service{
		Commands: Commands{
			CreateRestaurant: commands.NewCreateRestaurantHandler(restaurantRepo),
		},
		Queries: Queries{
			GetRestaurant: queries.NewGetRestaurantHandler(restaurantRepo),
		},
	}
}

func (s *Service) CreateRestaurant(ctx context.Context, cmd commands.CreateRestaurant) (string, error) {
	return s.Commands.CreateRestaurant.Handle(ctx, cmd)
}

func (s *Service) GetRestaurant(ctx context.Context, query queries.GetRestaurant) (*domain.Restaurant, error) {
	return s.Queries.GetRestaurant.Handle(ctx, query)
}
