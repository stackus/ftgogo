package domain

import (
	"github.com/stackus/edat/core"
)

func registerConsumerSnapshots() {
	core.RegisterSnapshots(ConsumerSnapshot{})
}

type ConsumerSnapshot struct {
	Name string
}

func (ConsumerSnapshot) SnapshotName() string {
	return "domain-service.domain.ConsumerSnapshot"
}
