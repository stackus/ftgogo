package queries

import (
	"context"
	"time"

	"github.com/stackus/ftgogo/order-history/internal/domain"
	"serviceapis/orderapi"
)

type GetConsumerOrderHistory struct {
	ConsumerID string
	Filter     *OrderHistoryFilters
	Next       *string
	Limit      *int
}

type GetConsumerOrderHistoryHandler struct {
	repo domain.OrderHistoryRepository
}

type OrderHistoryFilters struct {
	Since    *time.Time           `json:"since"`
	Keywords []*string            `json:"keywords"`
	Status   *orderapi.OrderState `json:"status"`
}

type GetConsumerOrderHistoryResponse struct {
	Orders []OrderHistory `json:"orders"`
	Next   string         `json:"next"`
}

func NewGetConsumerOrderHistoryHandler(orderHistoryRepo domain.OrderHistoryRepository) GetConsumerOrderHistoryHandler {
	return GetConsumerOrderHistoryHandler{repo: orderHistoryRepo}
}

func (h GetConsumerOrderHistoryHandler) Handle(ctx context.Context, query GetConsumerOrderHistory) (*GetConsumerOrderHistoryResponse, error) {
	filters := domain.OrderHistoryFilters{}

	if query.Next != nil {
		filters.Next = *query.Next
	}

	filters.Limit = domain.OrderHistoryLimit
	if query.Limit != nil {
		if *query.Limit >= 1 && *query.Limit <= domain.OrderHistoryMaximum {
			filters.Limit = *query.Limit
		}
	}

	if query.Filter != nil {
		for _, keyword := range query.Filter.Keywords {
			filters.Keywords = append(filters.Keywords, *keyword)
		}

		if query.Filter.Status != nil {
			filters.Status = *query.Filter.Status
		}

		if query.Filter.Since != nil {
			filters.Since = *query.Filter.Since
		}
	}

	orders, next, err := h.repo.FindConsumerOrders(ctx, query.ConsumerID, filters)
	if err != nil {
		return nil, err
	}

	history := make([]OrderHistory, 0, len(orders))
	for _, order := range orders {
		history = append(history, OrderHistory{
			OrderID:        order.OrderID,
			Status:         order.Status.String(),
			RestaurantID:   order.RestaurantID,
			RestaurantName: order.RestaurantName,
			CreatedAt:      order.CreatedAt,
		})
	}

	return &GetConsumerOrderHistoryResponse{
		Orders: history,
		Next:   next,
	}, nil
}
