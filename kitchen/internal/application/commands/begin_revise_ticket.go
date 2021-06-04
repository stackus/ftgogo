package commands

import (
	"context"

	"github.com/stackus/ftgogo/kitchen/internal/adapters"
	"github.com/stackus/ftgogo/kitchen/internal/domain"
)

type BeginReviseTicket struct {
	TicketID          string
	RestaurantID      string
	RevisedQuantities map[string]int
}

type BeginReviseTicketHandler struct {
	repo adapters.TicketRepository
}

func NewBeginReviseTicketHandler(ticketRepo adapters.TicketRepository) BeginReviseTicketHandler {
	return BeginReviseTicketHandler{repo: ticketRepo}
}

func (h BeginReviseTicketHandler) Handle(ctx context.Context, cmd BeginReviseTicket) error {
	_, err := h.repo.Update(ctx, cmd.TicketID, &domain.ReviseTicket{})

	return err
}
