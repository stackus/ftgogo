package steps

import (
	"encoding/json"
	"reflect"

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
	app     application.ServiceApplication
	account *domain.Account
	err     error
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
	f.err = nil

	accountRepo := adapters.NewAccountAggregateRepository(inmem.NewEventStore())
	f.app = application.NewServiceApplication(accountRepo)
}

func (f *FeatureState) RegisterCommonSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I expect the (?:request|command|query) to fail$`, f.iExpectTheCommandToFail)
	ctx.Step(`^I expect the (?:request|command|query) to succeed$`, f.iExpectTheCommandToSucceed)

	ctx.Step(`^(?:ensure )?the returned account matches:$`, f.theReturnedAccountMatches)
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

func (f *FeatureState) theReturnedAccountMatches(doc *godog.DocString) error {
	var err error
	var expected, actual AccountJson

	if err = json.Unmarshal([]byte(doc.Content), &expected); err != nil {
		return errors.Wrap(errors.ErrUnprocessableEntity, err.Error())
	}

	if f.account == nil {
		return errors.Wrap(errors.ErrNotFound, "Expected account to not be nil")
	}

	actual.ID = f.account.ID()
	actual.Name = f.account.Name
	actual.Enabled = f.account.Enabled
	if !reflect.DeepEqual(expected, actual) {
		return errors.Wrapf(errors.ErrInvalidArgument, "does not match expected: %v: got: %v", expected, actual)
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
