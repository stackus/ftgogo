package commands

import (
	"context"

	"github.com/stackus/ftgogo/delivery/internal/application/ports"
)

type SetCourierAvailability struct {
	CourierID string
	Available bool
}

type SetCourierAvailabilityHandler struct {
	repo ports.CourierRepository
}

func NewSetCourierAvailabilityHandler(courierRepo ports.CourierRepository) SetCourierAvailabilityHandler {
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
