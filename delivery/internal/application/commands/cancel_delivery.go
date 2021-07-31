package commands

import (
	"context"

	"github.com/stackus/ftgogo/delivery/internal/application/ports"
)

type CancelDelivery struct {
	OrderID string
}

type CancelDeliveryHandler struct {
	deliveryRepo ports.DeliveryRepository
	courierRepo  ports.CourierRepository
}

func NewCancelDeliveryHandler(deliveryRepo ports.DeliveryRepository, courierRepo ports.CourierRepository) CancelDeliveryHandler {
	return CancelDeliveryHandler{
		deliveryRepo: deliveryRepo,
		courierRepo:  courierRepo,
	}
}

func (h CancelDeliveryHandler) Handle(ctx context.Context, cmd CancelDelivery) error {
	delivery, err := h.deliveryRepo.Find(ctx, cmd.OrderID)
	if err != nil {
		return err
	}

	if delivery.AssignedCourierID != "" {
		courier, err := h.courierRepo.Find(ctx, delivery.AssignedCourierID)
		if err != nil {
			return err
		}

		courier.CancelDelivery(delivery.DeliveryID)

		err = h.courierRepo.Update(ctx, courier.CourierID, courier)
		if err != nil {
			return err
		}
	}

	delivery.Cancel()

	return h.deliveryRepo.Update(ctx, delivery.DeliveryID, delivery)
}
