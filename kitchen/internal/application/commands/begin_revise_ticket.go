package commands

import (
	"context"

	"github.com/stackus/ftgogo/kitchen/internal/application/ports"
	"github.com/stackus/ftgogo/kitchen/internal/domain"
)

type BeginReviseTicket struct {
	TicketID string
}

type BeginReviseTicketHandler struct {
	repo ports.TicketRepository
}

func NewBeginReviseTicketHandler(ticketRepo ports.TicketRepository) BeginReviseTicketHandler {
	return BeginReviseTicketHandler{repo: ticketRepo}
}

func (h BeginReviseTicketHandler) Handle(ctx context.Context, cmd BeginReviseTicket) error {
	_, err := h.repo.Update(ctx, cmd.TicketID, &domain.ReviseTicket{})

	return err
}
