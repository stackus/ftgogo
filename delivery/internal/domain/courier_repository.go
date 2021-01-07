package domain

import (
	"context"
)

type CourierRepository interface {
	Find(ctx context.Context, courierID string) (*Courier, error)
	FindOrCreate(ctx context.Context, courierID string) (*Courier, error)
	FindFirstAvailable(ctx context.Context) (*Courier, error)
	Save(ctx context.Context, courier *Courier) error
	Update(ctx context.Context, courierID string, courier *Courier) error
}
