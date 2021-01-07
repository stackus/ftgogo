package commands

import (
	"context"
	"strings"
	"time"

	"github.com/stackus/ftgogo/order-history/internal/domain"
	"serviceapis/orderapi"
)

type CreateOrderHistory struct {
	OrderID        string
	ConsumerID     string
	RestaurantID   string
	RestaurantName string
	LineItems      []orderapi.LineItem
	OrderTotal     int
}

type CreateOrderHistoryHandler struct {
	repo domain.OrderHistoryRepository
}

func NewCreateOrderHistoryHandler(orderHistoryRepo domain.OrderHistoryRepository) CreateOrderHistoryHandler {
	return CreateOrderHistoryHandler{repo: orderHistoryRepo}
}

func (h CreateOrderHistoryHandler) Handle(ctx context.Context, cmd CreateOrderHistory) error {
	keywords := []string{}
	seenWords := map[string]struct{}{}

	for _, word := range strings.Split(strings.ToLower(cmd.RestaurantName), " ") {
		if _, seen := seenWords[word]; !seen {
			keywords = append(keywords, word)
			seenWords[word] = struct{}{}
		}
	}

	for _, item := range cmd.LineItems {
		for _, word := range strings.Split(strings.ToLower(item.Name), " ") {
			if _, seen := seenWords[word]; !seen {
				keywords = append(keywords, word)
				seenWords[word] = struct{}{}
			}
		}
	}

	return h.repo.Save(ctx, &domain.OrderHistory{
		OrderID:        cmd.OrderID,
		ConsumerID:     cmd.ConsumerID,
		RestaurantID:   cmd.RestaurantID,
		RestaurantName: cmd.RestaurantName,
		LineItems:      cmd.LineItems,
		OrderTotal:     cmd.OrderTotal,
		Status:         orderapi.ApprovalPending,
		Keywords:       keywords,
		CreatedAt:      time.Now(),
	})
}
