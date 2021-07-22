package steps

import (
	"context"

	"github.com/cucumber/godog"
	"github.com/stackus/errors"

	"github.com/stackus/ftgogo/accounting/internal/application/queries"
)

func (f *FeatureState) RegisterGetAccountSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I request the account for "([^"]*)"$`, f.iRequestTheAccountFor)

	ctx.Step(`^(?:ensure )?the returned account is disabled$`, f.theReturnedAccountIsDisabled)
	ctx.Step(`^(?:ensure )?the returned account is enabled$`, f.theReturnedAccountIsEnabled)
}

func (f *FeatureState) iRequestTheAccountFor(consumerName string) error {
	accountID := f.accountNames[consumerName]

	query := queries.GetAccount{AccountID: accountID}

	f.account, f.err = f.app.GetAccount(context.Background(), query)

	return nil
}

func (f *FeatureState) theReturnedAccountIsDisabled() error {
	if f.account == nil {
		return errors.Wrap(errors.ErrNotFound, "expected account to not be nil")
	}

	if f.account.Enabled {
		return errors.Wrap(errors.ErrInvalidArgument, "expected the account to be disabled")
	}

	return nil
}

func (f *FeatureState) theReturnedAccountIsEnabled() error {
	if f.account == nil {
		return errors.Wrap(errors.ErrNotFound, "expected account to not be nil")
	}

	if !f.account.Enabled {
		return errors.Wrap(errors.ErrInvalidArgument, "expected the account to be enabled")
	}

	return nil
}
