package steps

import (
	"context"

	"github.com/cucumber/godog"
	"github.com/google/uuid"

	"github.com/stackus/ftgogo/delivery/internal/application/commands"
	"github.com/stackus/ftgogo/delivery/internal/application/queries"
)

func (f *FeatureState) RegisterSetCourierAvailabilitySteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^a courier exists named "([^"]*)"$`, f.aCourierExistsNamed)

	ctx.Step(`^I set the courier "([^"]*)" to be (available|unavailable)$`, f.iSetTheCourierToBe)

	ctx.Step(`^I get the courier named "([^"]*)"$`, f.iGetTheCourierNamed)
}

func (f *FeatureState) aCourierExistsNamed(courierName string) error {
	courierID := f.courierIDs[courierName]
	if courierID == "" {
		courierID = uuid.New().String()
		f.courierIDs[courierName] = courierID
	}

	f.err = f.app.SetCourierAvailability(context.Background(), commands.SetCourierAvailability{
		CourierID: courierID,
		Available: true,
	})

	return nil
}

func (f *FeatureState) iSetTheCourierToBe(courierName, availability string) error {
	courierID := f.courierIDs[courierName]

	available := true
	if availability == "unavailable" {
		available = false
	}

	f.err = f.app.SetCourierAvailability(context.Background(), commands.SetCourierAvailability{
		CourierID: courierID,
		Available: available,
	})

	return nil
}

func (f *FeatureState) iGetTheCourierNamed(courierName string) error {
	courierID := f.courierIDs[courierName]

	f.courier, f.err = f.app.GetCourier(context.Background(), queries.GetCourier{CourierID: courierID})

	return nil
}
