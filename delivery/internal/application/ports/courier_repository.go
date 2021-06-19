package ports

import (
	"context"

	"github.com/stackus/ftgogo/delivery/internal/domain"
)

type CourierRepository interface {
	Find(ctx context.Context, courierID string) (*domain.Courier, error)
	FindOrCreate(ctx context.Context, courierID string) (*domain.Courier, error)
	FindFirstAvailable(ctx context.Context) (*domain.Courier, error)
	Save(ctx context.Context, courier *domain.Courier) error
	Update(ctx context.Context, courierID string, courier *domain.Courier) error
}
