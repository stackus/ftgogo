package kitchenapi

import (
	"github.com/stackus/edat/core"
)

func registerReplies() {
	core.RegisterReplies(CreateTicketReply{})
}

type CreateTicketReply struct {
	TicketID string
}

func (CreateTicketReply) ReplyName() string { return "kitchenapi.CreateTicketReply" }
