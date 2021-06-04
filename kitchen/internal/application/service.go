package application

import (
	"github.com/stackus/ftgogo/kitchen/internal/application/commands"
	"github.com/stackus/ftgogo/kitchen/internal/application/queries"
)

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
	GetRestaurant queries.GetRestaurantHandler
}
