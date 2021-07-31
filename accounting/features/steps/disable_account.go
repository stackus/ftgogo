package steps

import (
	"context"

	"github.com/cucumber/godog"

	"github.com/stackus/ftgogo/accounting/internal/application/commands"
)

func (f *FeatureState) RegisterDisableAccountSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I disable the account for "([^"]*)"$`, f.iDisableTheAccountFor)
}

func (f *FeatureState) iDisableTheAccountFor(consumerName string) error {
	accountID := f.accountNames[consumerName]

	cmd := commands.DisableAccount{AccountID: accountID}

	f.err = f.app.DisableAccount(context.Background(), cmd)

	return nil
}
