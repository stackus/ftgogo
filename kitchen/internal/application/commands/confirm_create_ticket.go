package commands

import (
	"context"

	"github.com/stackus/ftgogo/kitchen/internal/domain"
)

type ConfirmCreateTicket struct {
	TicketID string
}

type ConfirmCreateTicketHandler struct {
	repo      domain.TicketRepository
	publisher domain.TicketPublisher
}

func NewConfirmCreateTicketHandler(ticketRepo domain.TicketRepository, ticketPublisher domain.TicketPublisher) ConfirmCreateTicketHandler {
	return ConfirmCreateTicketHandler{
		repo:      ticketRepo,
		publisher: ticketPublisher,
	}
}

func (h ConfirmCreateTicketHandler) Handle(ctx context.Context, cmd ConfirmCreateTicket) error {
	root, err := h.repo.Update(ctx, cmd.TicketID, &domain.ConfirmCreateTicket{})
	if err != nil {
		return err
	}

	return h.publisher.PublishEntityEvents(ctx, root)
}
