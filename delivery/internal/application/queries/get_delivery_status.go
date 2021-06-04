package queries

import (
	"context"

	"github.com/stackus/ftgogo/delivery/internal/adapters"
)

type GetDeliveryStatus struct {
	DeliveryID string
}

type GetDeliveryStatusHandler struct {
	deliveryRepo adapters.DeliveryRepository
	courierRepo  adapters.CourierRepository
}

type DeliveryStatus struct {
	ID              string   `json:"id"`
	AssignedCourier string   `json:"assigned_courier"`
	CourierActions  []string `json:"courier_actions"`
	Status          string   `json:"status"`
}

func NewGetDeliveryStatusHandler(deliveryRepo adapters.DeliveryRepository, courierRepo adapters.CourierRepository) GetDeliveryStatusHandler {
	return GetDeliveryStatusHandler{deliveryRepo: deliveryRepo}
}

func (h GetDeliveryStatusHandler) Handle(ctx context.Context, query GetDeliveryStatus) (*DeliveryStatus, error) {
	delivery, err := h.deliveryRepo.Find(ctx, query.DeliveryID)
	if err != nil {
		return nil, err
	}

	deliveryStatus := &DeliveryStatus{
		ID:     delivery.DeliveryID,
		Status: delivery.Status.String(),
	}

	if delivery.AssignedCourierID != "" {
		courier, err := h.courierRepo.Find(ctx, delivery.AssignedCourierID)
		if err != nil {
			return nil, err
		}

		deliveryStatus.AssignedCourier = courier.CourierID

		for _, action := range courier.Plan.ActionsFor(delivery.DeliveryID) {
			deliveryStatus.CourierActions = append(deliveryStatus.CourierActions, action.ActionType.String())
		}
	}

	return deliveryStatus, nil
}
