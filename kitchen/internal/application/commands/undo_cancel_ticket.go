package commands

import (
	"context"

	"github.com/stackus/ftgogo/kitchen/internal/application/ports"
	"github.com/stackus/ftgogo/kitchen/internal/domain"
)

type UndoCancelTicket struct {
	TicketID     string
	RestaurantID string
}

type UndoCancelTicketHandler struct {
	repo ports.TicketRepository
}

func NewUndoCancelTicketHandler(ticketRepo ports.TicketRepository) UndoCancelTicketHandler {
	return UndoCancelTicketHandler{repo: ticketRepo}
}

func (h UndoCancelTicketHandler) Handle(ctx context.Context, cmd UndoCancelTicket) error {
	_, err := h.repo.Update(ctx, cmd.TicketID, &domain.UndoCancelTicket{})

	return err
}
