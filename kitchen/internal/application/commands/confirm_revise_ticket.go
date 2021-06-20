package commands

import (
	"context"

	"github.com/stackus/ftgogo/kitchen/internal/application/ports"
	"github.com/stackus/ftgogo/kitchen/internal/domain"
)

type ConfirmReviseTicket struct {
	TicketID          string
	RestaurantID      string
	RevisedQuantities map[string]int
}

type ConfirmReviseTicketHandler struct {
	repo ports.TicketRepository
}

func NewConfirmReviseTicketHandler(repo ports.TicketRepository) ConfirmReviseTicketHandler {
	return ConfirmReviseTicketHandler{
		repo: repo,
	}
}

func (h ConfirmReviseTicketHandler) Handle(ctx context.Context, cmd ConfirmReviseTicket) error {
	_, err := h.repo.Update(ctx, cmd.TicketID, &domain.ConfirmReviseTicket{})
	if err != nil {
		return err
	}

	return nil
}
