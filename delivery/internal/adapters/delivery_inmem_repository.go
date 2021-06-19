package adapters

import (
	"context"

	"github.com/stackus/errors"

	"github.com/stackus/ftgogo/delivery/internal/application/ports"
	"github.com/stackus/ftgogo/delivery/internal/domain"
)

type DeliveryInmemRepository struct {
	deliveries map[string]*domain.Delivery
}

var _ ports.DeliveryRepository = (*DeliveryInmemRepository)(nil)

func NewDeliveryInmemRepository() *DeliveryInmemRepository {
	return &DeliveryInmemRepository{deliveries: map[string]*domain.Delivery{}}
}

func (r *DeliveryInmemRepository) Find(ctx context.Context, deliveryID string) (*domain.Delivery, error) {
	if delivery, exists := r.deliveries[deliveryID]; !exists {
		return nil, domain.ErrDeliveryNotFound
	} else {
		return delivery, nil
	}
}

func (r *DeliveryInmemRepository) Save(ctx context.Context, delivery *domain.Delivery) error {
	if _, exists := r.deliveries[delivery.DeliveryID]; exists {
		return errors.Wrap(errors.ErrConflict, "delivery already exists")
	}
	r.deliveries[delivery.DeliveryID] = delivery
	return nil
}

func (r *DeliveryInmemRepository) Update(ctx context.Context, deliveryID string, delivery *domain.Delivery) error {
	if _, exists := r.deliveries[deliveryID]; !exists {
		return domain.ErrDeliveryNotFound
	}
	r.deliveries[deliveryID] = delivery
	return nil
}
