package commands

import (
	"context"

	"github.com/stackus/ftgogo/kitchen/internal/application/ports"
	"github.com/stackus/ftgogo/kitchen/internal/domain"
)

type CancelCreateTicket struct {
	TicketID string
}

type CancelCreateTicketHandler struct {
	repo ports.TicketRepository
}

func NewCancelCreateTicketHandler(ticketRepo ports.TicketRepository) CancelCreateTicketHandler {
	return CancelCreateTicketHandler{
		repo: ticketRepo,
	}
}

func (h CancelCreateTicketHandler) Handle(ctx context.Context, cmd CancelCreateTicket) error {
	_, err := h.repo.Update(ctx, cmd.TicketID, &domain.CancelCreateTicket{})

	return err
}
