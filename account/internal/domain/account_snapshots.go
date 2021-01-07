package domain

import (
	"github.com/stackus/edat/core"
)

func registerAccountSnapshots() {
	core.RegisterSnapshots(AccountSnapshot{})
}

type AccountSnapshot struct {
	Name    string
	Enabled bool
}

func (AccountSnapshot) SnapshotName() string { return "accountingservice.AccountSnapshot" }
