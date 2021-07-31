package adapters

import (
	"context"

	"github.com/stackus/ftgogo/store-web/internal/domain"
)

type (
	FindDelivery struct {
		DeliveryID string
	}

	FindCourier struct {
		CourierID string
	}

	SetCourierAvailability struct {
		CourierID string
		Available bool
	}
)

type DeliveryRepository interface {
	FindDelivery(context.Context, FindDelivery) (*domain.Delivery, error)
	FindCourier(context.Context, FindCourier) (*domain.Courier, error)
	SetCourierAvailability(context.Context, SetCourierAvailability) error
}
