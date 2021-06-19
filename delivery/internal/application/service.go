package application

import (
	"context"

	"github.com/stackus/ftgogo/delivery/internal/application/commands"
	"github.com/stackus/ftgogo/delivery/internal/application/ports"
	"github.com/stackus/ftgogo/delivery/internal/application/queries"
	"github.com/stackus/ftgogo/delivery/internal/domain"
)

type ServiceApplication interface {
	CreateDelivery(ctx context.Context, cmd commands.CreateDelivery) error
	CreateRestaurant(ctx context.Context, cmd commands.CreateRestaurant) error
	SetCourierAvailability(ctx context.Context, cmd commands.SetCourierAvailability) error
	ScheduleDelivery(ctx context.Context, cmd commands.ScheduleDelivery) error
	CancelDelivery(ctx context.Context, cmd commands.CancelDelivery) error
	GetCourier(ctx context.Context, query queries.GetCourier) (*domain.Courier, error)
	GetDelivery(ctx context.Context, query queries.GetDelivery) (*domain.Delivery, error)
}

type Service struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateDelivery         commands.CreateDeliveryHandler
	CreateRestaurant       commands.CreateRestaurantHandler
	SetCourierAvailability commands.SetCourierAvailabilityHandler
	ScheduleDelivery       commands.ScheduleDeliveryHandler
	CancelDelivery         commands.CancelDeliveryHandler
}

type Queries struct {
	GetCourier  queries.GetCourierHandler
	GetDelivery queries.GetDeliveryHandler
}

func NewServiceApplication(courierRepo ports.CourierRepository, deliveryRepo ports.DeliveryRepository, restaurantRepo ports.RestaurantRepository) *Service {
	return &Service{
		Commands: Commands{
			CreateDelivery:         commands.NewCreateDeliveryHandler(deliveryRepo, restaurantRepo),
			CreateRestaurant:       commands.NewCreateRestaurantHandler(restaurantRepo),
			SetCourierAvailability: commands.NewSetCourierAvailabilityHandler(courierRepo),
			ScheduleDelivery:       commands.NewScheduleDeliveryHandler(deliveryRepo, courierRepo),
			CancelDelivery:         commands.NewCancelDeliveryHandler(deliveryRepo, courierRepo),
		},
		Queries: Queries{
			GetCourier:  queries.NewGetCourierHandler(courierRepo),
			GetDelivery: queries.NewGetDeliveryHandler(deliveryRepo),
		},
	}
}

func (s Service) CreateDelivery(ctx context.Context, cmd commands.CreateDelivery) error {
	return s.Commands.CreateDelivery.Handle(ctx, cmd)
}

func (s Service) CreateRestaurant(ctx context.Context, cmd commands.CreateRestaurant) error {
	return s.Commands.CreateRestaurant.Handle(ctx, cmd)
}

func (s Service) SetCourierAvailability(ctx context.Context, cmd commands.SetCourierAvailability) error {
	return s.Commands.SetCourierAvailability.Handle(ctx, cmd)
}

func (s Service) ScheduleDelivery(ctx context.Context, cmd commands.ScheduleDelivery) error {
	return s.Commands.ScheduleDelivery.Handle(ctx, cmd)
}

func (s Service) CancelDelivery(ctx context.Context, cmd commands.CancelDelivery) error {
	return s.Commands.CancelDelivery.Handle(ctx, cmd)
}

func (s Service) GetCourier(ctx context.Context, query queries.GetCourier) (*domain.Courier, error) {
	return s.Queries.GetCourier.Handle(ctx, query)
}

func (s Service) GetDelivery(ctx context.Context, query queries.GetDelivery) (*domain.Delivery, error) {
	return s.Queries.GetDelivery.Handle(ctx, query)
}
