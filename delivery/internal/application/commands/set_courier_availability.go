package commands

import (
	"context"

	"github.com/stackus/ftgogo/delivery/internal/adapters"
)

type SetCourierAvailability struct {
	CourierID string
	Available bool
}

type SetCourierAvailabilityHandler struct {
	repo adapters.CourierRepository
}

func NewSetCourierAvailabilityHandler(courierRepo adapters.CourierRepository) SetCourierAvailabilityHandler {
	return SetCourierAvailabilityHandler{repo: courierRepo}
}

func (h SetCourierAvailabilityHandler) Handle(ctx context.Context, cmd SetCourierAvailability) error {
	courier, err := h.repo.FindOrCreate(ctx, cmd.CourierID)
	if err != nil {
		return err
	}

	courier.Available = cmd.Available

	return h.repo.Update(ctx, cmd.CourierID, courier)
}
