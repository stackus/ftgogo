package steps

import (
	"context"

	"github.com/cucumber/godog"

	"github.com/stackus/ftgogo/delivery/internal/application/commands"
)

func (f *FeatureState) RegisterCancelDeliverySteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I cancel delivery for order "([^"]*)"$`, f.iCancelDeliveryForOrder)
}

func (f *FeatureState) iCancelDeliveryForOrder(orderID string) error {
	f.err = f.app.CancelDelivery(context.Background(), commands.CancelDelivery{OrderID: orderID})

	return nil
}
