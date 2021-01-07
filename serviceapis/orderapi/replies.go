package orderapi

import (
	"github.com/stackus/edat/core"
)

func registerReplies() {
	core.RegisterReplies(BeginReviseOrderReply{})
}

type BeginReviseOrderReply struct {
	RevisedOrderTotal int
}

func (BeginReviseOrderReply) ReplyName() string { return "orderapi.BeginReviseOrderReply" }
