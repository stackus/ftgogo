package commands

import (
	"context"
	"time"

	"github.com/stackus/ftgogo/kitchen/internal/domain"
)

type AcceptTicket struct {
	TicketID string
	ReadyBy  time.Time
}

type AcceptTicketHandler struct {
	repo      domain.TicketRepository
	publisher domain.TicketPublisher
}

func NewAcceptTicketHandler(ticketRepo domain.TicketRepository, ticketPublisher domain.TicketPublisher) AcceptTicketHandler {
	return AcceptTicketHandler{
		repo:      ticketRepo,
		publisher: ticketPublisher,
	}
}

func (h AcceptTicketHandler) Handle(ctx context.Context, cmd AcceptTicket) error {
	root, err := h.repo.Update(ctx, cmd.TicketID, &domain.AcceptTicket{ReadyBy: cmd.ReadyBy})
	if err != nil {
		return err
	}

	return h.publisher.PublishEntityEvents(ctx, root)
}
