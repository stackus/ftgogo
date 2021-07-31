package queries

import (
	"context"

	"github.com/stackus/ftgogo/delivery/internal/application/ports"
	"github.com/stackus/ftgogo/delivery/internal/domain"
)

type GetDelivery struct {
	OrderID string
}

type GetDeliveryHandler struct {
	repo ports.DeliveryRepository
}

func NewGetDeliveryHandler(repo ports.DeliveryRepository) GetDeliveryHandler {
	return GetDeliveryHandler{
		repo: repo,
	}
}

func (h GetDeliveryHandler) Handle(ctx context.Context, query GetDelivery) (*domain.Delivery, error) {
	return h.repo.Find(ctx, query.OrderID)
}
