package domain

import (
	"github.com/stackus/edat/core"

	"github.com/stackus/ftgogo/serviceapis/commonapi"
)

func registerConsumerSnapshots() {
	core.RegisterSnapshots(ConsumerSnapshot{})
}

type ConsumerSnapshot struct {
	Name      string
	Addresses map[string]*commonapi.Address
}

func (ConsumerSnapshot) SnapshotName() string {
	return "consumerservice.ConsumerSnapshot"
}
