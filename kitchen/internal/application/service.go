package application

import (
	"context"

	"github.com/stackus/ftgogo/kitchen/internal/application/commands"
	"github.com/stackus/ftgogo/kitchen/internal/application/ports"
	"github.com/stackus/ftgogo/kitchen/internal/application/queries"
	"github.com/stackus/ftgogo/kitchen/internal/domain"
)

type ServiceApplication interface {
	CreateTicket(ctx context.Context, cmd commands.CreateTicket) (string, error)
	ConfirmCreateTicket(ctx context.Context, cmd commands.ConfirmCreateTicket) error
	CancelCreateTicket(ctx context.Context, cmd commands.CancelCreateTicket) error
	BeginCancelTicket(ctx context.Context, cmd commands.BeginCancelTicket) error
	UndoCancelTicket(ctx context.Context, cmd commands.UndoCancelTicket) error
	ConfirmCancelTicket(ctx context.Context, cmd commands.ConfirmCancelTicket) error
	BeginReviseTicket(ctx context.Context, cmd commands.BeginReviseTicket) error
	UndoReviseTicket(ctx context.Context, cmd commands.UndoReviseTicket) error
	ConfirmReviseTicket(ctx context.Context, cmd commands.ConfirmReviseTicket) error
	AcceptTicket(ctx context.Context, cmd commands.AcceptTicket) error
	CreateRestaurant(ctx context.Context, cmd commands.CreateRestaurant) error
	ReviseRestaurantMenu(ctx context.Context, cmd commands.ReviseRestaurantMenu) error
	GetTicket(ctx context.Context, query queries.GetTicket) (*domain.Ticket, error)
	GetRestaurant(ctx context.Context, query queries.GetRestaurant) (*domain.Restaurant, error)
}

type Service struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateTicket         commands.CreateTicketHandler
	ConfirmCreateTicket  commands.ConfirmCreateTicketHandler
	CancelCreateTicket   commands.CancelCreateTicketHandler
	BeginCancelTicket    commands.BeginCancelTicketHandler
	ConfirmCancelTicket  commands.ConfirmCancelTicketHandler
	UndoCancelTicket     commands.UndoCancelTicketHandler
	BeginReviseTicket    commands.BeginReviseTicketHandler
	ConfirmReviseTicket  commands.ConfirmReviseTicketHandler
	UndoReviseTicket     commands.UndoReviseTicketHandler
	AcceptTicket         commands.AcceptTicketHandler
	CreateRestaurant     commands.CreateRestaurantHandler
	ReviseRestaurantMenu commands.ReviseRestaurantMenuHandler
}

type Queries struct {
	GetTicket     queries.GetTicketHandler
	GetRestaurant queries.GetRestaurantHandler
}

func NewServiceApplication(ticketRepo ports.TicketRepository, restaurantRepo ports.RestaurantRepository) *Service {
	return &Service{
		Commands: Commands{
			CreateTicket:         commands.NewCreateTicketHandler(ticketRepo),
			ConfirmCreateTicket:  commands.NewConfirmCreateTicketHandler(ticketRepo),
			CancelCreateTicket:   commands.NewCancelCreateTicketHandler(ticketRepo),
			BeginCancelTicket:    commands.NewBeginCancelTicketHandler(ticketRepo),
			ConfirmCancelTicket:  commands.NewConfirmCancelTicketHandler(ticketRepo),
			UndoCancelTicket:     commands.NewUndoCancelTicketHandler(ticketRepo),
			BeginReviseTicket:    commands.NewBeginReviseTicketHandler(ticketRepo),
			ConfirmReviseTicket:  commands.NewConfirmReviseTicketHandler(ticketRepo),
			UndoReviseTicket:     commands.NewUndoReviseTicketHandler(ticketRepo),
			AcceptTicket:         commands.NewAcceptTicketHandler(ticketRepo),
			CreateRestaurant:     commands.NewCreateRestaurantHandler(restaurantRepo),
			ReviseRestaurantMenu: commands.NewReviseRestaurantMenuHandler(restaurantRepo),
		},
		Queries: Queries{
			GetTicket:     queries.NewGetTicketHandler(ticketRepo),
			GetRestaurant: queries.NewGetRestaurantHandler(restaurantRepo),
		},
	}
}

func (s Service) CreateTicket(ctx context.Context, cmd commands.CreateTicket) (string, error) {
	return s.Commands.CreateTicket.Handle(ctx, cmd)
}

func (s Service) ConfirmCreateTicket(ctx context.Context, cmd commands.ConfirmCreateTicket) error {
	return s.Commands.ConfirmCreateTicket.Handle(ctx, cmd)
}

func (s Service) CancelCreateTicket(ctx context.Context, cmd commands.CancelCreateTicket) error {
	return s.Commands.CancelCreateTicket.Handle(ctx, cmd)
}

func (s Service) BeginCancelTicket(ctx context.Context, cmd commands.BeginCancelTicket) error {
	return s.Commands.BeginCancelTicket.Handle(ctx, cmd)
}

func (s Service) UndoCancelTicket(ctx context.Context, cmd commands.UndoCancelTicket) error {
	return s.Commands.UndoCancelTicket.Handle(ctx, cmd)
}

func (s Service) ConfirmCancelTicket(ctx context.Context, cmd commands.ConfirmCancelTicket) error {
	return s.Commands.ConfirmCancelTicket.Handle(ctx, cmd)
}

func (s Service) BeginReviseTicket(ctx context.Context, cmd commands.BeginReviseTicket) error {
	return s.Commands.BeginReviseTicket.Handle(ctx, cmd)
}

func (s Service) UndoReviseTicket(ctx context.Context, cmd commands.UndoReviseTicket) error {
	return s.Commands.UndoReviseTicket.Handle(ctx, cmd)
}

func (s Service) ConfirmReviseTicket(ctx context.Context, cmd commands.ConfirmReviseTicket) error {
	return s.Commands.ConfirmReviseTicket.Handle(ctx, cmd)
}

func (s Service) AcceptTicket(ctx context.Context, cmd commands.AcceptTicket) error {
	return s.Commands.AcceptTicket.Handle(ctx, cmd)
}

func (s Service) CreateRestaurant(ctx context.Context, cmd commands.CreateRestaurant) error {
	return s.Commands.CreateRestaurant.Handle(ctx, cmd)
}

func (s Service) ReviseRestaurantMenu(ctx context.Context, cmd commands.ReviseRestaurantMenu) error {
	return s.Commands.ReviseRestaurantMenu.Handle(ctx, cmd)
}

func (s Service) GetTicket(ctx context.Context, query queries.GetTicket) (*domain.Ticket, error) {
	return s.Queries.GetTicket.Handle(ctx, query)
}

func (s Service) GetRestaurant(ctx context.Context, query queries.GetRestaurant) (*domain.Restaurant, error) {
	return s.Queries.GetRestaurant.Handle(ctx, query)
}
