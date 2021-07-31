package ports

import (
	"context"

	"github.com/stackus/ftgogo/delivery/internal/domain"
)

type DeliveryRepository interface {
	Find(ctx context.Context, deliveryID string) (*domain.Delivery, error)
	Save(ctx context.Context, delivery *domain.Delivery) error
	Update(ctx context.Context, deliveryID string, delivery *domain.Delivery) error
}
