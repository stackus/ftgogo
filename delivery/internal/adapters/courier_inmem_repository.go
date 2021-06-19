package adapters

import (
	"context"

	"github.com/google/uuid"
	"github.com/stackus/errors"

	"github.com/stackus/ftgogo/delivery/internal/application/ports"
	"github.com/stackus/ftgogo/delivery/internal/domain"
)

type CourierInmemRepository struct {
	couriers map[string]*domain.Courier
}

var _ ports.CourierRepository = (*CourierInmemRepository)(nil)

func NewCourierInmemRepository() *CourierInmemRepository {
	return &CourierInmemRepository{couriers: map[string]*domain.Courier{}}
}

func (r *CourierInmemRepository) Find(ctx context.Context, courierID string) (*domain.Courier, error) {
	if courier, exists := r.couriers[courierID]; !exists {
		return nil, domain.ErrCourierNotFound
	} else {
		return courier, nil
	}
}

func (r *CourierInmemRepository) FindOrCreate(ctx context.Context, courierID string) (*domain.Courier, error) {
	if courier, exists := r.couriers[courierID]; !exists {
		courier = &domain.Courier{
			CourierID: courierID,
			Plan:      domain.Plan{},
			Available: true,
		}
		r.couriers[courierID] = courier
		return courier, nil
	} else {
		return courier, nil
	}
}

func (r *CourierInmemRepository) FindFirstAvailable(ctx context.Context) (*domain.Courier, error) {
	if len(r.couriers) > 0 {
		for _, courier := range r.couriers {
			if courier.Available {
				return courier, nil
			}
		}
	}

	return r.FindOrCreate(ctx, uuid.New().String())
}

func (r *CourierInmemRepository) Save(ctx context.Context, courier *domain.Courier) error {
	if _, exists := r.couriers[courier.CourierID]; exists {
		return errors.Wrap(errors.ErrConflict, "courier with that identifier already exists")
	}
	r.couriers[courier.CourierID] = courier
	return nil
}

func (r *CourierInmemRepository) Update(ctx context.Context, courierID string, courier *domain.Courier) error {
	if _, exists := r.couriers[courierID]; !exists {
		return domain.ErrCourierNotFound
	}
	r.couriers[courierID] = courier
	return nil
}
