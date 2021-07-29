package steps

import (
	"context"

	"github.com/cucumber/godog"

	"github.com/stackus/ftgogo/order/internal/application/commands"
)

func (f *FeatureState) RegisterReviseOrderSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I (?:have )?beg(?:i|u)n to revise the order to "([^"]*)" from "([^"]*)"$`, f.iBeginToReviseTheOrderToFrom)
	ctx.Step(`^I (?:have )?confirm(?:ed)? revising the order to "([^"]*)" from "([^"]*)"$`, f.iConfirmRevisingTheOrderToFrom)
	ctx.Step(`^I (?:have )?undo(?:ne)? revising the order to "([^"]*)" from "([^"]*)"$`, f.iUndoRevisingTheOrderToFrom)

}

func (f *FeatureState) iBeginToReviseTheOrderToFrom(restaurantName, consumerName string) error {
	orderID := f.orderIDs[restaurantName+consumerName]

	_, f.err = f.app.BeginReviseOrder(context.Background(), commands.BeginReviseOrder{OrderID: orderID})

	return nil
}

func (f *FeatureState) iConfirmRevisingTheOrderToFrom(restaurantName, consumerName string) error {
	orderID := f.orderIDs[restaurantName+consumerName]

	f.err = f.app.ConfirmReviseOrder(context.Background(), commands.ConfirmReviseOrder{OrderID: orderID})

	return nil
}

func (f *FeatureState) iUndoRevisingTheOrderToFrom(restaurantName, consumerName string) error {
	orderID := f.orderIDs[restaurantName+consumerName]

	f.err = f.app.UndoReviseOrder(context.Background(), commands.UndoReviseOrder{OrderID: orderID})

	return nil
}
