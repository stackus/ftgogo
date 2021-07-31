package steps

import (
	"context"
	"time"

	"github.com/cucumber/godog"

	"github.com/stackus/ftgogo/order/internal/application/commands"
)

func (f *FeatureState) RegisterCreateOrderSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I (?:have )?submit(?:ted)? an order to "([^"]*)" from "([^"]*)"$`, f.iSubmitAnOrderToFrom)
	ctx.Step(`^I (?:have )?approved? the order to "([^"]*)" from "([^"]*)" with ticket "([^"]*)"$`, f.iApproveTheOrderToFromWithTicket)
	ctx.Step(`^I (?:have )?reject(?:ed)? the order to "([^"]*)" from "([^"]*)"$`, f.iRejectTheOrderToFrom)
}

func (f *FeatureState) iSubmitAnOrderToFrom(restaurantName, consumerName string) error {
	restaurant, err := getRestaurantFromFixture(restaurantName)
	if err != nil {
		return err
	}

	consumer, err := getConsumerFromFixture(consumerName)
	if err != nil {
		return err
	}

	f.orderIDs[restaurantName+consumerName], f.err = f.app.CreateOrder(context.Background(), commands.CreateOrder{
		ConsumerID:   consumer.ID,
		RestaurantID: restaurant.RestaurantID,
		DeliverAt:    time.Now(),
		DeliverTo:    consumer.Address,
		LineItems:    nil,
	})

	return nil
}

func (f *FeatureState) iApproveTheOrderToFromWithTicket(restaurantName, consumerName, ticketID string) error {
	orderID := f.orderIDs[restaurantName+consumerName]

	f.err = f.app.ApproveOrder(context.Background(), commands.ApproveOrder{
		OrderID:  orderID,
		TicketID: ticketID,
	})

	if f.err == nil {
		f.ticketIDs[restaurantName+consumerName] = ticketID
	}

	return nil
}

func (f *FeatureState) iRejectTheOrderToFrom(restaurantName, consumerName string) error {
	orderID := f.orderIDs[restaurantName+consumerName]

	f.err = f.app.RejectOrder(context.Background(), commands.RejectOrder{
		OrderID: orderID,
	})

	return nil
}
