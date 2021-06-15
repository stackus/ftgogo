package commands

import (
	"context"

	"github.com/stackus/ftgogo/store-web/internal/adapters"
)

type SetCourierAvailability struct {
	CourierID string
	Available bool
}

type SetCourierAvailabilityHandler struct {
	repo adapters.DeliveryRepository
}

func NewSetCourierAvailabilityHandler(repo adapters.DeliveryRepository) SetCourierAvailabilityHandler {
	return SetCourierAvailabilityHandler{repo: repo}
}

func (h SetCourierAvailabilityHandler) Handle(ctx context.Context, cmd SetCourierAvailability) error {
	return h.repo.SetCourierAvailability(ctx, adapters.SetCourierAvailability{
		CourierID: cmd.CourierID,
		Available: cmd.Available,
	})
}
