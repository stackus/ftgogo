package domain

import (
	"time"

	"serviceapis/commonapi"
	"shared-go/errs"
)

type ActionType string

const (
	PickUp  ActionType = "PICKUP"
	DropOff            = "DROPOFF"
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
	Address    commonapi.Address
	When       time.Time
}

// Courier errors
var (
	ErrCourierNotFound = errs.NewError("courier not found", errs.ErrNotFound)
)

func (a ActionType) String() string {
	return string(a)
}

func (c *Courier) AddAction(action Action) {
	c.Plan.Add(action)
}

func (c *Courier) CancelDelivery(deliveryID string) {
	c.Plan.RemoveDelivery(deliveryID)
}

func (p *Plan) Add(action Action) {
	*p = append(*p, action)
}

func (p *Plan) RemoveDelivery(deliveryID string) {
	replacement := Plan{}
	for _, action := range *p {
		if !action.IsFor(deliveryID) {
			replacement = append(replacement, action)
		}
	}

	*p = replacement
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
