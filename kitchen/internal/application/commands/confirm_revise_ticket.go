package commands

import (
	"context"

	"github.com/stackus/ftgogo/kitchen/internal/domain"
)

type ConfirmReviseTicket struct {
	TicketID          string
	RestaurantID      string
	RevisedQuantities map[string]int
}

type ConfirmReviseTicketHandler struct {
	repo      domain.TicketRepository
	publisher domain.TicketPublisher
}

func NewConfirmReviseTicketHandler(ticketRepo domain.TicketRepository, ticketPublisher domain.TicketPublisher) ConfirmReviseTicketHandler {
	return ConfirmReviseTicketHandler{
		repo:      ticketRepo,
		publisher: ticketPublisher,
	}
}

func (h ConfirmReviseTicketHandler) Handle(ctx context.Context, cmd ConfirmReviseTicket) error {
	root, err := h.repo.Update(ctx, cmd.TicketID, &domain.ConfirmReviseTicket{})
	if err != nil {
		return err
	}

	return h.publisher.PublishEntityEvents(ctx, root)
}
