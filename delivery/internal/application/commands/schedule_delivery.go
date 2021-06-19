package commands

import (
	"context"
	"time"

	"github.com/stackus/ftgogo/delivery/internal/application/ports"
	"github.com/stackus/ftgogo/delivery/internal/domain"
)

type ScheduleDelivery struct {
	OrderID string
	ReadyBy time.Time
}

type ScheduleDeliveryHandler struct {
	deliveryRepo ports.DeliveryRepository
	courierRepo  ports.CourierRepository
}

func NewScheduleDeliveryHandler(deliveryRepo ports.DeliveryRepository, courierRepo ports.CourierRepository) ScheduleDeliveryHandler {
	return ScheduleDeliveryHandler{
		deliveryRepo: deliveryRepo,
		courierRepo:  courierRepo,
	}
}

func (h ScheduleDeliveryHandler) Handle(ctx context.Context, cmd ScheduleDelivery) error {
	delivery, err := h.deliveryRepo.Find(ctx, cmd.OrderID)
	if err != nil {
		return err
	}

	courier, err := h.courierRepo.FindFirstAvailable(ctx)
	if err != nil {
		return err
	}

	courier.Plan.Add(domain.Action{
		DeliveryID: delivery.DeliveryID,
		ActionType: domain.PickUp,
		Address:    delivery.PickUpAddress,
		When:       cmd.ReadyBy,
	})
	courier.Plan.Add(domain.Action{
		DeliveryID: delivery.DeliveryID,
		ActionType: domain.DropOff,
		Address:    delivery.DeliveryAddress,
		When:       cmd.ReadyBy.Add(30 * time.Minute),
	})

	err = h.courierRepo.Update(ctx, courier.CourierID, courier)
	if err != nil {
		return err
	}

	delivery.Schedule(cmd.ReadyBy, courier.CourierID)

	return h.deliveryRepo.Update(ctx, delivery.DeliveryID, delivery)
}
