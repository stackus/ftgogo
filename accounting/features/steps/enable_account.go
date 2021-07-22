package steps

import (
	"context"

	"github.com/cucumber/godog"

	"github.com/stackus/ftgogo/accounting/internal/application/commands"
)

func (f *FeatureState) RegisterEnableAccountSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I enable the account for "([^"]*)"$`, f.iEnableTheAccountFor)
}

func (f *FeatureState) iEnableTheAccountFor(consumerName string) error {
	accountID := f.accountNames[consumerName]

	cmd := commands.EnableAccount{AccountID: accountID}

	f.err = f.app.EnableAccount(context.Background(), cmd)

	return nil
}
