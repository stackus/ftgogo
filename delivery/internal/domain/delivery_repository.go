package domain

import (
	"context"
)

type DeliveryRepository interface {
	Find(ctx context.Context, deliveryID string) (*Delivery, error)
	Save(ctx context.Context, delivery *Delivery) error
	Update(ctx context.Context, deliveryID string, delivery *Delivery) error
}
