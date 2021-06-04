package commands

import (
	"context"

	"github.com/stackus/ftgogo/kitchen/internal/adapters"
	"github.com/stackus/ftgogo/kitchen/internal/domain"
)

type UndoReviseTicket struct {
	TicketID     string
	RestaurantID string
}

type UndoReviseTicketHandler struct {
	repo adapters.TicketRepository
}

func NewUndoReviseTicketHandler(ticketRepo adapters.TicketRepository) UndoReviseTicketHandler {
	return UndoReviseTicketHandler{repo: ticketRepo}
}

func (h UndoReviseTicketHandler) Handle(ctx context.Context, cmd UndoReviseTicket) error {
	_, err := h.repo.Update(ctx, cmd.TicketID, &domain.UndoReviseTicket{})

	return err
}
