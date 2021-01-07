package consumerapi

import (
	"github.com/stackus/edat/core"
)

func registerReplies() {
	core.RegisterReplies(ConsumerNotFound{})
}

type ConsumerNotFound struct{}

func (ConsumerNotFound) ReplyName() string { return "consumerapi.ConsumerNotFound" }
