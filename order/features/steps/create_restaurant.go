package steps

import (
	"context"

	"github.com/cucumber/godog"

	"github.com/stackus/ftgogo/order/internal/application/commands"
)

func (f *FeatureState) RegisterCreateRestaurantSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I have (?:created|initialized) the restaurant "([^"]*)"$`, f.iHaveInitializedTheRestaurant)
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
