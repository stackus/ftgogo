package accountingapi

import (
	"github.com/stackus/edat/core"
)

func registerReplies() {
	core.RegisterReplies(AccountDisabled{})
}

type AccountDisabled struct{}

func (AccountDisabled) ReplyName() string { return "accountingapi.AccountDisabled" }
