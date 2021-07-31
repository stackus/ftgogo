package steps

import (
	"context"

	"github.com/cucumber/godog"
	"github.com/google/uuid"

	"github.com/stackus/ftgogo/delivery/internal/application/commands"
)

func (f *FeatureState) RegisterCreateRestaurantSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^(?:I create )?(?:a|another) restaurant named "([^"]*)" (?:exists )?with address$`, f.aRestaurantNamedExistsWithAddress)
}

func (f *FeatureState) aRestaurantNamedExistsWithAddress(restaurantName string, table *godog.Table) error {
	address, err := parseAddressFromTable(table)
	if err != nil {
		return err
	}

	restaurantID := f.restaurantIDs[restaurantName]
	if restaurantID == "" {
		restaurantID = uuid.New().String()
		f.restaurantIDs[restaurantName] = restaurantID
	}

	f.err = f.app.CreateRestaurant(context.Background(), commands.CreateRestaurant{
		RestaurantID: restaurantID,
		Name:         restaurantName,
		Address:      address,
	})

	return nil
}
