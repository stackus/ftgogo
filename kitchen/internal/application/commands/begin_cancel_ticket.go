package commands

import (
	"context"

	"github.com/stackus/ftgogo/kitchen/internal/domain"
)

type BeginCancelTicket struct {
	TicketID     string
	RestaurantID string
}

type BeginCancelTicketHandler struct {
	repo domain.TicketRepository
}

func NewBeginCancelTicketHandler(ticketRepo domain.TicketRepository) BeginCancelTicketHandler {
	return BeginCancelTicketHandler{repo: ticketRepo}
}

func (h BeginCancelTicketHandler) Handle(ctx context.Context, cmd BeginCancelTicket) error {
	_, err := h.repo.Update(ctx, cmd.TicketID, &domain.CancelTicket{})

	return err
}
