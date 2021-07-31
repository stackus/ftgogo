package queries

import (
	"context"

	"github.com/stackus/ftgogo/store-web/internal/adapters"
	"github.com/stackus/ftgogo/store-web/internal/domain"
)

type GetDeliveryHistory struct {
	DeliveryID string
}

type GetDeliveryHistoryHandler struct {
	repo adapters.DeliveryRepository
}

func NewGetDeliveryHistoryHandler(repo adapters.DeliveryRepository) GetDeliveryHistoryHandler {
	return GetDeliveryHistoryHandler{
		repo: repo,
	}
}

func (h GetDeliveryHistoryHandler) Handle(ctx context.Context, query GetDeliveryHistory) (*domain.DeliveryHistory, error) {
	delivery, err := h.repo.FindDelivery(ctx, adapters.FindDelivery{DeliveryID: query.DeliveryID})
	if err != nil {
		return nil, err
	}

	history := &domain.DeliveryHistory{
		ID:     delivery.DeliveryID,
		Status: delivery.Status.String(),
	}

	if delivery.AssignedCourierID != "" {
		courier, err := h.repo.FindCourier(ctx, adapters.FindCourier{CourierID: delivery.AssignedCourierID})
		if err != nil {
			return nil, err
		}

		history.AssignedCourier = courier.CourierID

		for _, action := range courier.Plan.ActionsFor(delivery.DeliveryID) {
			history.CourierActions = append(history.CourierActions, action.ActionType.String())
		}
	}

	return history, nil
}
