package commands

import (
	"context"

	"github.com/stackus/ftgogo/kitchen/internal/adapters"
	"github.com/stackus/ftgogo/kitchen/internal/domain"
	"github.com/stackus/ftgogo/serviceapis/kitchenapi"
)

type CreateTicket struct {
	OrderID      string
	RestaurantID string
	LineItems    []kitchenapi.LineItem
}

type CreateTicketHandler struct {
	repo adapters.TicketRepository
}

func NewCreateTicketHandler(ticketRepo adapters.TicketRepository) CreateTicketHandler {
	return CreateTicketHandler{repo: ticketRepo}
}

func (h CreateTicketHandler) Handle(ctx context.Context, cmd CreateTicket) (string, error) {
	root, err := h.repo.Save(ctx, &domain.CreateTicket{
		OrderID:      cmd.OrderID,
		RestaurantID: cmd.RestaurantID,
		LineItems:    cmd.LineItems,
	})
	if err != nil {
		return "", err
	}

	return root.AggregateID(), nil
}
