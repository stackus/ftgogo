package steps

import (
	"context"
	"time"

	"github.com/cucumber/godog"

	"github.com/stackus/ftgogo/delivery/internal/application/commands"
)

func (f *FeatureState) RegisterScheduleDeliverySteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I schedule the delivery for order "([^"]*)"$`, f.iScheduleTheDeliveryForOrder)
}

func (f *FeatureState) iScheduleTheDeliveryForOrder(orderID string) error {
	f.err = f.app.ScheduleDelivery(context.Background(), commands.ScheduleDelivery{
		OrderID: orderID,
		ReadyBy: time.Now(),
	})

	return nil
}
