package domain

import (
	"github.com/stackus/edat/core"
)

func registerReviseOrderSagaData() {
	core.RegisterSagaData(ReviseOrderSagaData{})
}

type ReviseOrderSagaData struct {
	OrderID           string
	ConsumerID        string
	RestaurantID      string
	TicketID          string
	ExpectedVersion   int
	RevisedOrderTotal int
	RevisedQuantities map[string]int
}

func (ReviseOrderSagaData) SagaDataName() string {
	return "orderservice.ReviseOrderSagaData"
}
