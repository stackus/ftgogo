package queries

import (
	"context"

	"github.com/stackus/ftgogo/delivery/internal/adapters"
	"github.com/stackus/ftgogo/delivery/internal/domain"
)

type GetDelivery struct {
	DeliveryID string
}

type GetDeliveryHandler struct {
	repo adapters.DeliveryRepository
}

func NewGetDeliveryHandler(repo adapters.DeliveryRepository) GetDeliveryHandler {
	return GetDeliveryHandler{
		repo: repo,
	}
}

func (h GetDeliveryHandler) Handle(ctx context.Context, query GetDelivery) (*domain.Delivery, error) {
	return h.repo.Find(ctx, query.DeliveryID)
}
