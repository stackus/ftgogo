package commands

import (
	"context"

	"github.com/stackus/ftgogo/kitchen/internal/domain"
	"serviceapis/commonapi"
)

type BeginReviseTicket struct {
	TicketID          string
	RestaurantID      string
	RevisedQuantities commonapi.MenuItemQuantities
}

type BeginReviseTicketHandler struct {
	repo domain.TicketRepository
}

func NewBeginReviseTicketHandler(ticketRepo domain.TicketRepository) BeginReviseTicketHandler {
	return BeginReviseTicketHandler{repo: ticketRepo}
}

func (h BeginReviseTicketHandler) Handle(ctx context.Context, cmd BeginReviseTicket) error {
	_, err := h.repo.Update(ctx, cmd.TicketID, &domain.ReviseTicket{})

	return err
}
