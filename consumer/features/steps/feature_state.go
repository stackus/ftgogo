package steps

import (
	"encoding/json"
	"reflect"

	"github.com/cucumber/godog"
	_ "github.com/stackus/edat-msgpack"
	"github.com/stackus/edat/inmem"
	"github.com/stackus/errors"

	"github.com/stackus/ftgogo/consumer/internal/adapters"
	"github.com/stackus/ftgogo/consumer/internal/application"
	"github.com/stackus/ftgogo/consumer/internal/domain"
	"github.com/stackus/ftgogo/serviceapis"
	"github.com/stackus/ftgogo/serviceapis/commonapi"
)

type ConsumerJson struct {
	ID        string
	Name      string
	Addresses map[string]*commonapi.Address
}

type FeatureState struct {
	app        application.ServiceApplication
	consumerID string
	consumer   *domain.Consumer
	address    *commonapi.Address
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
	f.consumerID = ""
	f.consumer = nil
	f.address = nil
	f.err = nil

	accountRepo := adapters.NewConsumerAggregateRepository(inmem.NewEventStore())
	f.app = application.NewServiceApplication(accountRepo)
}

func (f *FeatureState) RegisterCommonSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I expect the (?:request|command|query) to fail$`, f.iExpectTheCommandToFail)
	ctx.Step(`^I expect the (?:request|command|query) to succeed$`, f.iExpectTheCommandToSucceed)

	ctx.Step(`^(?:ensure )?the returned consumer matches:$`, f.theReturnedConsumerMatches)
	ctx.Step(`^(?:ensure )?the returned address matches:$`, f.theReturnedAddressMatches)
	ctx.Step(`^(?:ensure )?the returned error message is:$`, f.theReturnedErrorMessageIs)
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

func (f *FeatureState) theReturnedConsumerMatches(doc *godog.DocString) error {
	var err error
	var expected, actual ConsumerJson

	if err = json.Unmarshal([]byte(doc.Content), &expected); err != nil {
		return errors.Wrap(errors.ErrUnprocessableEntity, err.Error())
	}

	if f.consumer == nil {
		return errors.Wrap(errors.ErrNotFound, "Expected consumer to not be nil")
	}

	actual.ID = f.consumer.ID()
	actual.Name = f.consumer.Name()
	actual.Addresses = f.consumer.Addresses()
	if !reflect.DeepEqual(expected, actual) {
		return errors.Wrapf(errors.ErrInvalidArgument, "does not match expected: %v: got: %v", expected, actual)
	}

	return nil
}

func (f *FeatureState) theReturnedAddressMatches(doc *godog.DocString) error {
	var err error
	var expected *commonapi.Address

	if err = json.Unmarshal([]byte(doc.Content), &expected); err != nil {
		return errors.Wrap(errors.ErrUnprocessableEntity, err.Error())
	}

	if f.address == nil {
		return errors.Wrap(errors.ErrNotFound, "Expected address to not be nil")
	}

	if !reflect.DeepEqual(expected, f.address) {
		return errors.Wrapf(errors.ErrInvalidArgument, "does not match expected: %v: got: %v", expected, f.address)
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
