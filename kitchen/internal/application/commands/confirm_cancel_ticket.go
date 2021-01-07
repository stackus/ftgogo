package commands

import (
	"context"

	"github.com/stackus/ftgogo/kitchen/internal/domain"
)

type ConfirmCancelTicket struct {
	TicketID     string
	RestaurantID string
}

type ConfirmCancelTicketHandler struct {
	repo      domain.TicketRepository
	publisher domain.TicketPublisher
}

func NewConfirmCancelTicketHandler(ticketRepo domain.TicketRepository, ticketPublisher domain.TicketPublisher) ConfirmCancelTicketHandler {
	return ConfirmCancelTicketHandler{
		repo:      ticketRepo,
		publisher: ticketPublisher,
	}
}

func (h ConfirmCancelTicketHandler) Handle(ctx context.Context, cmd ConfirmCancelTicket) error {
	root, err := h.repo.Update(ctx, cmd.TicketID, &domain.ConfirmCancelTicket{})
	if err != nil {
		return err
	}

	return h.publisher.PublishEntityEvents(ctx, root)
}
