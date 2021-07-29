package steps

import (
	"context"

	"github.com/cucumber/godog"

	"github.com/stackus/ftgogo/order/internal/application/commands"
)

func (f *FeatureState) RegisterCancelOrderSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I (?:have )?beg(?:i|u)n to cancel the order to "([^"]*)" from "([^"]*)"$`, f.iBeginToCancelTheOrderToFrom)
	ctx.Step(`^I (?:have )?confirm(?:ed)? cancelling the order to "([^"]*)" from "([^"]*)"$`, f.iConfirmCancellingTheOrderToFrom)
	ctx.Step(`^I (?:have )?undo(?:ne)? cancelling the order to "([^"]*)" from "([^"]*)"$`, f.iUndoCancellingTheOrderToFrom)

}

func (f *FeatureState) iBeginToCancelTheOrderToFrom(restaurantName, consumerName string) error {
	orderID := f.orderIDs[restaurantName+consumerName]

	f.err = f.app.BeginCancelOrder(context.Background(), commands.BeginCancelOrder{OrderID: orderID})

	return nil
}

func (f *FeatureState) iConfirmCancellingTheOrderToFrom(restaurantName, consumerName string) error {
	orderID := f.orderIDs[restaurantName+consumerName]

	f.err = f.app.ConfirmCancelOrder(context.Background(), commands.ConfirmCancelOrder{OrderID: orderID})

	return nil
}

func (f *FeatureState) iUndoCancellingTheOrderToFrom(restaurantName, consumerName string) error {
	orderID := f.orderIDs[restaurantName+consumerName]

	f.err = f.app.UndoCancelOrder(context.Background(), commands.UndoCancelOrder{OrderID: orderID})

	return nil
}
