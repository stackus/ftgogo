package steps

import (
	"context"

	"github.com/cucumber/godog"
	"github.com/stackus/errors"

	"github.com/stackus/ftgogo/order/internal/application/queries"
)

func (f *FeatureState) RegisterGetOrderSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^(?:ensure |expect )?the order to "([^"]*)" from "([^"]*)" is "([^"]*)"$`, f.theOrderToFromIs)
}

func (f *FeatureState) theOrderToFromIs(restaurantName, consumerName, expected string) error {
	orderID := f.orderIDs[restaurantName+consumerName]

	order, err := f.app.GetOrder(context.Background(), queries.GetOrder{OrderID: orderID})
	if err != nil {
		return err
	}

	got := order.State.String()
	if got != expected {
		return errors.Wrapf(errors.ErrInvalidArgument, "order state does not match expected: %s: got: %s", expected, got)
	}

	return nil
}
