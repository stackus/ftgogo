package queries

import (
	"context"

	"github.com/stackus/ftgogo/delivery/internal/adapters"
	"github.com/stackus/ftgogo/delivery/internal/domain"
)

type GetCourier struct {
	CourierID string
}

type GetCourierHandler struct {
	repo adapters.CourierRepository
}

func NewGetCourierHandler(repo adapters.CourierRepository) GetCourierHandler {
	return GetCourierHandler{repo: repo}
}

func (h GetCourierHandler) Handle(ctx context.Context, query GetCourier) (*domain.Courier, error) {
	return h.repo.Find(ctx, query.CourierID)
}
