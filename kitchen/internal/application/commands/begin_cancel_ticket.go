package commands

import (
	"context"

	"github.com/stackus/ftgogo/kitchen/internal/application/ports"
	"github.com/stackus/ftgogo/kitchen/internal/domain"
)

type BeginCancelTicket struct {
	TicketID string
}

type BeginCancelTicketHandler struct {
	repo ports.TicketRepository
}

func NewBeginCancelTicketHandler(ticketRepo ports.TicketRepository) BeginCancelTicketHandler {
	return BeginCancelTicketHandler{repo: ticketRepo}
}

func (h BeginCancelTicketHandler) Handle(ctx context.Context, cmd BeginCancelTicket) error {
	_, err := h.repo.Update(ctx, cmd.TicketID, &domain.CancelTicket{})

	return err
}
