package commands

import (
	"context"

	"github.com/stackus/ftgogo/kitchen/internal/adapters"
	"github.com/stackus/ftgogo/kitchen/internal/domain"
)

type ConfirmCancelTicket struct {
	TicketID     string
	RestaurantID string
}

type ConfirmCancelTicketHandler struct {
	repo adapters.TicketRepository
}

func NewConfirmCancelTicketHandler(repo adapters.TicketRepository) ConfirmCancelTicketHandler {
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
