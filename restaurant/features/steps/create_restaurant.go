package steps

import (
	"context"

	"github.com/cucumber/godog"

	"github.com/stackus/ftgogo/restaurant/internal/application/commands"
)

func (f *FeatureState) RegisterCreateRestaurantSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I (?:have )?(?:create|initialize)d? the restaurant "([^"]*)"$`, f.iCreateTheRestaurant)
}

func (f *FeatureState) iCreateTheRestaurant(restaurantName string) error {
	restaurant, err := getRestaurantFromFixture(restaurantName)
	if err != nil {
		return err
	}

	f.restaurantIDs[restaurantName], f.err = f.app.CreateRestaurant(context.Background(), commands.CreateRestaurant{
		Name:      restaurant.Name,
		Address:   restaurant.Address,
		MenuItems: restaurant.MenuItems,
	})

	return nil
}
