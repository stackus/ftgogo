package domain

import (
	"time"

	"github.com/stackus/ftgogo/serviceapis/commonapi"
)

type ActionType string

const (
	PickUp  ActionType = "PICKUP"
	DropOff ActionType = "DROPOFF"
)

type Courier struct {
	CourierID string
	Plan      Plan
	Available bool
}

type Plan []Action

type Action struct {
	DeliveryID string
	ActionType ActionType
	Address    *commonapi.Address
	When       time.Time
}

func (a ActionType) String() string {
	return string(a)
}

func (p Plan) ActionsFor(deliveryID string) []Action {
	actions := []Action{}
	for _, action := range p {
		if action.IsFor(deliveryID) {
			actions = append(actions, action)
		}
	}

	return actions
}

func (a Action) IsFor(deliveryID string) bool {
	return a.DeliveryID == deliveryID
}
