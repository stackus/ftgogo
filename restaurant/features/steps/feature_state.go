package steps

import (
	"github.com/cucumber/godog"
	"github.com/rdumont/assistdog"
	_ "github.com/stackus/edat-msgpack"
	"github.com/stackus/errors"

	"github.com/stackus/ftgogo/restaurant/internal/adapters"
	"github.com/stackus/ftgogo/restaurant/internal/application"
	"github.com/stackus/ftgogo/serviceapis"
)

type FeatureState struct {
	app           application.ServiceApplication
	restaurantIDs map[string]string
	err           error
}

func NewFeatureState() *FeatureState {
	f := &FeatureState{}
	f.Reset()

	return f
}

var assist = assistdog.NewDefault()

func init() {
	serviceapis.RegisterTypes()
}

func (f *FeatureState) Reset() {
	f.restaurantIDs = make(map[string]string)
	f.err = nil

	restaurantRepo := adapters.NewRestaurantInmemRepository()
	f.app = application.NewServiceApplication(restaurantRepo)
}

func (f *FeatureState) RegisterCommonSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I expect the (?:request|command|query) to fail$`, f.iExpectTheCommandToFail)
	ctx.Step(`^I expect the (?:request|command|query) to succeed$`, f.iExpectTheCommandToSucceed)

	ctx.Step(`^(?:ensure |expect )?the returned error message is "([^"]*)"$`, f.theReturnedErrorMessageIs)
}

func (f *FeatureState) iExpectTheCommandToFail() error {
	if f.err == nil {
		return errors.Wrap(errors.ErrUnknown, "expected error to not be nil")
	}
	return nil
}

func (f *FeatureState) iExpectTheCommandToSucceed() error {
	if f.err != nil {
		return errors.Wrapf(f.err, "expected error to be nil: got %w", f.err)
	}

	return nil
}

func (f *FeatureState) theReturnedErrorMessageIs(errorMsg string) error {
	if f.err == nil {
		return errors.Wrap(errors.ErrUnknown, "expected error to not be nil")
	}

	if errorMsg != f.err.Error() {
		return errors.Wrapf(errors.ErrInvalidArgument, "expected: %s: got: %s", errorMsg, f.err.Error())
	}

	return nil
}
