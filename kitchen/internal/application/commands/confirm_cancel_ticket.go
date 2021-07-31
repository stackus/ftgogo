package commands

import (
	"context"

	"github.com/stackus/ftgogo/kitchen/internal/application/ports"
	"github.com/stackus/ftgogo/kitchen/internal/domain"
)

type ConfirmCancelTicket struct {
	TicketID     string
	RestaurantID string
}

type ConfirmCancelTicketHandler struct {
	repo ports.TicketRepository
}

func NewConfirmCancelTicketHandler(repo ports.TicketRepository) ConfirmCancelTicketHandler {
	return ConfirmCancelTicketHandler{
		repo: repo,
	}
}

func (h ConfirmCancelTicketHandler) Handle(ctx context.Context, cmd ConfirmCancelTicket) error {
	_, err := h.repo.Update(ctx, cmd.TicketID, &domain.ConfirmCancelTicket{})
	if err != nil {
		return err
	}

	return nil
}
