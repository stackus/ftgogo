package steps

import (
	"context"
	"encoding/json"

	"github.com/cucumber/godog"

	"github.com/stackus/ftgogo/kitchen/internal/application/commands"
)

func (f *FeatureState) RegisterCreateRestaurantSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I (?:create|setup) (?:a|the|another) restaurant with:$`, f.iCreateARestaurantWith)
}

func (f *FeatureState) iCreateARestaurantWith(doc *godog.DocString) error {
	var cmd commands.CreateRestaurant

	err := json.Unmarshal([]byte(doc.Content), &cmd)
	if err != nil {
		return err
	}

	f.err = f.app.CreateRestaurant(context.Background(), cmd)

	return nil
}
