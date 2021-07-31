package steps

import (
	"github.com/cucumber/godog"
	"github.com/rdumont/assistdog"
	_ "github.com/stackus/edat-msgpack"
	"github.com/stackus/edat/inmem"
	"github.com/stackus/errors"

	"github.com/stackus/ftgogo/kitchen/internal/adapters"
	"github.com/stackus/ftgogo/kitchen/internal/application"
	"github.com/stackus/ftgogo/kitchen/internal/domain"
	"github.com/stackus/ftgogo/serviceapis"
)

type FeatureState struct {
	app           application.ServiceApplication
	ticket        *domain.Ticket
	ticketID      string
	restaurant    *domain.Restaurant
	restaurantIDs map[string]string
	ticketIDs     map[string]string
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
	domain.RegisterTypes()
}

func (f *FeatureState) Reset() {
	f.ticket = nil
	f.restaurant = nil
	f.restaurantIDs = make(map[string]string)
	f.ticketIDs = make(map[string]string)
	f.err = nil

	ticketRepo := adapters.NewTicketAggregateRepository(inmem.NewEventStore())
	restaurantRepo := adapters.NewRestaurantInmemRepository()
	f.app = application.NewServiceApplication(ticketRepo, restaurantRepo)
}

func (f *FeatureState) RegisterCommonSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I expect the (?:request|command|query) to fail$`, f.iExpectTheCommandToFail)
	ctx.Step(`^I expect the (?:request|command|query) to succeed$`, f.iExpectTheCommandToSucceed)

	ctx.Step(`^(?:ensure |expect )?the returned error message is "([^"]*)"$`, f.theReturnedErrorMessageIs)

}

func (f *FeatureState) iExpectTheCommandToFail() error {
	if f.err == nil {
		return errors.Wrap(errors.ErrUnknown, "Expected error to not be nil")
	}
	return nil
}

func (f *FeatureState) iExpectTheCommandToSucceed() error {
	if f.err != nil {
		return errors.Wrap(f.err, "Expected error to be nil")
	}
	return nil
}

func (f *FeatureState) theReturnedErrorMessageIs(errorMsg string) error {
	if f.err == nil {
		return errors.Wrap(errors.ErrUnknown, "Expected error to not be nil")
	}

	if errorMsg != f.err.Error() {
		return errors.Wrapf(errors.ErrInvalidArgument, "expected: %s: got: %s", errorMsg, f.err.Error())
	}

	return nil
}
