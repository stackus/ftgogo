package queries

import (
	"context"

	"github.com/stackus/ftgogo/kitchen/internal/application/ports"
	"github.com/stackus/ftgogo/kitchen/internal/domain"
)

type GetTicket struct {
	TicketID string
}

type GetTicketHandler struct {
	repo ports.TicketRepository
}

func NewGetTicketHandler(repo ports.TicketRepository) GetTicketHandler {
	return GetTicketHandler{repo: repo}
}

func (h GetTicketHandler) Handle(ctx context.Context, query GetTicket) (*domain.Ticket, error) {
	return h.repo.Load(ctx, query.TicketID)
}
