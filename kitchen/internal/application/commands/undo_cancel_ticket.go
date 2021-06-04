package commands

import (
	"context"

	"github.com/stackus/ftgogo/kitchen/internal/adapters"
	"github.com/stackus/ftgogo/kitchen/internal/domain"
)

type UndoCancelTicket struct {
	TicketID     string
	RestaurantID string
}

type UndoCancelTicketHandler struct {
	repo adapters.TicketRepository
}

func NewUndoCancelTicketHandler(ticketRepo adapters.TicketRepository) UndoCancelTicketHandler {
	return UndoCancelTicketHandler{repo: ticketRepo}
}

func (h UndoCancelTicketHandler) Handle(ctx context.Context, cmd UndoCancelTicket) error {
	_, err := h.repo.Update(ctx, cmd.TicketID, &domain.UndoCancelTicket{})

	return err
}
