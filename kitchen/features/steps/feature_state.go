package steps

import (
	"encoding/json"
	"reflect"

	"github.com/cucumber/godog"
	_ "github.com/stackus/edat-msgpack"
	"github.com/stackus/edat/inmem"
	"github.com/stackus/errors"

	"github.com/stackus/ftgogo/kitchen/internal/adapters"
	"github.com/stackus/ftgogo/kitchen/internal/application"
	"github.com/stackus/ftgogo/kitchen/internal/domain"
	"github.com/stackus/ftgogo/serviceapis"
	"github.com/stackus/ftgogo/serviceapis/commonapi"
)

type FeatureState struct {
	app        application.ServiceApplication
	ticket     *domain.Ticket
	ticketID   string
	restaurant *domain.Restaurant
	err        error
}

func NewFeatureState() *FeatureState {
	f := &FeatureState{}
	f.Reset()

	return f
}

func init() {
	serviceapis.RegisterTypes()
	domain.RegisterTypes()
}

func (f *FeatureState) Reset() {
	f.ticket = nil
	f.restaurant = nil
	f.err = nil

	ticketRepo := adapters.NewTicketAggregateRepository(inmem.NewEventStore())
	restaurantRepo := adapters.NewRestaurantInmemRepository()
	f.app = application.NewServiceApplication(ticketRepo, restaurantRepo)
}

func (f *FeatureState) RegisterCommonSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I expect the (?:request|command|query) to fail$`, f.iExpectTheCommandToFail)
	ctx.Step(`^I expect the (?:request|command|query) to succeed$`, f.iExpectTheCommandToSucceed)

	ctx.Step(`^(?:ensure )?the returned error message is:$`, f.theReturnedErrorMessageIs)
	ctx.Step(`^(?:ensure )?the returned ticket matches:$`, f.theReturnedTicketMatches)
	ctx.Step(`^(?:ensure )?the returned ticket status is:$`, f.theReturnedTicketStatusIs)
	ctx.Step(`^(?:ensure )?the returned restaurant matches:$`, f.theReturnedRestaurantMatches)
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

func (f *FeatureState) theReturnedErrorMessageIs(doc *godog.DocString) error {
	if f.err == nil {
		return errors.Wrap(errors.ErrUnknown, "Expected error to not be nil")
	}

	if doc.Content != f.err.Error() {
		return errors.Wrapf(errors.ErrInvalidArgument, "expected: %s: got: %s", doc.Content, f.err.Error())
	}

	return nil
}

func (f *FeatureState) theReturnedTicketMatches(doc *godog.DocString) error {
	var err error
	var expected *domain.Ticket

	if err = json.Unmarshal([]byte(doc.Content), &expected); err != nil {
		return errors.Wrap(errors.ErrUnprocessableEntity, err.Error())
	}

	if f.ticket == nil {
		return errors.Wrap(errors.ErrNotFound, "Expected courier to not be nil")
	}

	if !reflect.DeepEqual(expected, f.ticket) {
		return errors.Wrapf(errors.ErrInvalidArgument, "does not match expected: %v: got: %v", expected, f.ticket)
	}

	return nil
}

func (f *FeatureState) theReturnedTicketStatusIs(doc *godog.DocString) error {
	if f.ticket == nil {
		return errors.Wrap(errors.ErrNotFound, "Expected courier to not be nil")
	}

	if f.ticket.State.String() != doc.Content {
		return errors.Wrapf(errors.ErrInvalidArgument, "does not match expected: %v: got: %v", doc.Content, f.ticket.State.String())
	}

	return nil
}

func (f *FeatureState) theReturnedRestaurantMatches(doc *godog.DocString) error {
	var err error
	var expected *commonapi.Address

	if err = json.Unmarshal([]byte(doc.Content), &expected); err != nil {
		return errors.Wrap(errors.ErrUnprocessableEntity, err.Error())
	}

	if f.restaurant == nil {
		return errors.Wrap(errors.ErrNotFound, "Expected delivery to not be nil")
	}

	if !reflect.DeepEqual(expected, f.restaurant) {
		return errors.Wrapf(errors.ErrInvalidArgument, "does not match expected: %v: got: %v", expected, f.restaurant)
	}

	return nil
}
