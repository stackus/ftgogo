package steps

import (
	"github.com/cucumber/godog"
	_ "github.com/stackus/edat-msgpack"
	"github.com/stackus/edat/inmem"
	"github.com/stackus/errors"

	"github.com/stackus/ftgogo/accounting/internal/adapters"
	"github.com/stackus/ftgogo/accounting/internal/application"
	"github.com/stackus/ftgogo/accounting/internal/domain"
	"github.com/stackus/ftgogo/serviceapis"
)

type AccountJson struct {
	ID      string
	Name    string
	Enabled bool
}

type FeatureState struct {
	app          application.ServiceApplication
	account      *domain.Account
	accountNames map[string]string
	err          error
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
	f.account = nil
	f.accountNames = make(map[string]string)
	f.err = nil

	accountRepo := adapters.NewAccountAggregateRepository(inmem.NewEventStore())
	f.app = application.NewServiceApplication(accountRepo)
}

func (f *FeatureState) RegisterCommonSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I expect the (?:request|command|query) to fail$`, f.iExpectTheCommandToFail)
	ctx.Step(`^I expect the (?:request|command|query) to succeed$`, f.iExpectTheCommandToSucceed)

	ctx.Step(`^the returned error message is "([^"]*)"$`, f.theReturnedErrorMessageIs)
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
