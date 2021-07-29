package steps

import (
	"context"

	"github.com/cucumber/godog"

	"github.com/stackus/ftgogo/order/internal/application/commands"
)

func (f *FeatureState) RegisterCreateRestaurantSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I have (?:created|initialized) the restaurant "([^"]*)"$`, f.iHaveInitializedTheRestaurant)
	// ctx.Step(`^I (?:create|setup|have created) (?:a|the|another) ticket for order "([^"]*)" (?:and|at) restaurant "([^"]*)" with items$`, f.iCreateATicketForOrderAndRestaurant)
	// ctx.Step(`^I (?:confirm|have confirmed) creating (?:a|the|another) ticket for order "([^"]*)"$`, f.iConfirmCreateATicketForOrder)
	// ctx.Step(`^I (?:cancel|have cancell?ed) creating (?:a|the|another) ticket for order "([^"]*)"$`, f.iCancelCreateATicketForOrder)
}

func (f *FeatureState) iHaveInitializedTheRestaurant(restaurantName string) error {
	restaurant, err := getRestaurantFromFixture(restaurantName)
	if err != nil {
		return err
	}

	f.err = f.app.CreateRestaurant(context.Background(), commands.CreateRestaurant{
		RestaurantID: restaurant.RestaurantID,
		Name:         restaurant.Name,
		Menu:         restaurant.MenuItems,
	})

	return nil
}
