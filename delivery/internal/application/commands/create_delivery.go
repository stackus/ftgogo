package commands

import (
	"context"
	"time"

	"github.com/stackus/ftgogo/delivery/internal/application/ports"
	"github.com/stackus/ftgogo/delivery/internal/domain"
	"github.com/stackus/ftgogo/serviceapis/commonapi"
)

type CreateDelivery struct {
	OrderID         string
	RestaurantID    string
	DeliveryAddress *commonapi.Address
}

type CreateDeliveryHandler struct {
	deliveryRepo   ports.DeliveryRepository
	restaurantRepo ports.RestaurantRepository
}

func NewCreateDeliveryHandler(deliveryRepo ports.DeliveryRepository, restaurantRepo ports.RestaurantRepository) CreateDeliveryHandler {
	return CreateDeliveryHandler{
		deliveryRepo:   deliveryRepo,
		restaurantRepo: restaurantRepo,
	}
}

func (h CreateDeliveryHandler) Handle(ctx context.Context, cmd CreateDelivery) error {
	restaurant, err := h.restaurantRepo.Find(ctx, cmd.RestaurantID)
	if err != nil {
		return err
	}

	delivery := &domain.Delivery{
		DeliveryID:        cmd.OrderID,
		RestaurantID:      restaurant.RestaurantID,
		AssignedCourierID: "",
		PickUpAddress:     restaurant.Address,
		DeliveryAddress:   cmd.DeliveryAddress,
		Status:            domain.DeliveryPending,
		PickUpTime:        time.Time{},
		ReadyBy:           time.Time{},
	}

	return h.deliveryRepo.Save(ctx, delivery)
}
