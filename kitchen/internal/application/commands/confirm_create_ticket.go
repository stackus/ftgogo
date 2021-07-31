package commands

import (
	"context"

	"github.com/stackus/ftgogo/kitchen/internal/application/ports"
	"github.com/stackus/ftgogo/kitchen/internal/domain"
)

type ConfirmCreateTicket struct {
	TicketID string
}

type ConfirmCreateTicketHandler struct {
	repo ports.TicketRepository
}

func NewConfirmCreateTicketHandler(repo ports.TicketRepository) ConfirmCreateTicketHandler {
	return ConfirmCreateTicketHandler{
		repo: repo,
	}
}

func (h ConfirmCreateTicketHandler) Handle(ctx context.Context, cmd ConfirmCreateTicket) error {
	_, err := h.repo.Update(ctx, cmd.TicketID, &domain.ConfirmCreateTicket{})
	if err != nil {
		return err
	}

	return nil
}
