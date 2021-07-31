package steps

import (
	"context"

	"github.com/cucumber/godog"

	"github.com/stackus/ftgogo/delivery/internal/application/commands"
)

func (f *FeatureState) RegisterCreateDeliverySteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I create (?:a|another) delivery for order "([^"]*)" from "([^"]*)" to address$`, f.iCreateADeliveryForOrderFromToAddress)
}

func (f *FeatureState) iCreateADeliveryForOrderFromToAddress(orderID, restaurantName string, table *godog.Table) error {
	address, err := parseAddressFromTable(table)
	if err != nil {
		return err
	}

	restaurantID := f.restaurantIDs[restaurantName]

	f.err = f.app.CreateDelivery(context.Background(), commands.CreateDelivery{
		OrderID:         orderID,
		RestaurantID:    restaurantID,
		DeliveryAddress: address,
	})

	return nil
}
