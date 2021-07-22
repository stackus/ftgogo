package steps

import (
	"github.com/cucumber/godog"
	"github.com/rdumont/assistdog"
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
	app                 application.ServiceApplication
	consumerID          string
	consumer            *domain.Consumer
	address             *commonapi.Address
	registeredConsumers map[string]string
	err                 error
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
	f.consumerID = ""
	f.consumer = nil
	f.address = nil
	f.registeredConsumers = make(map[string]string)
	f.err = nil

	accountRepo := adapters.NewConsumerAggregateRepository(inmem.NewEventStore())
	f.app = application.NewServiceApplication(accountRepo)
}

func (f *FeatureState) RegisterCommonSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I expect the (?:request|command|query) to fail$`, f.iExpectTheCommandToFail)
	ctx.Step(`^I expect the (?:request|command|query) to succeed$`, f.iExpectTheCommandToSucceed)

	ctx.Step(`^(?:ensure )?the returned error message is "([^"]*)"$`, f.theReturnedErrorMessageIs)
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
