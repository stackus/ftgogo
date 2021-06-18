package steps

import (
	"context"

	"github.com/cucumber/godog"

	"github.com/stackus/ftgogo/accounting/internal/application/queries"
)

func (f *FeatureState) RegisterGetAccountSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I (?:request|get|fetch) an account with:?$`, f.iRequestAnAccountWith)
}

func (f *FeatureState) iRequestAnAccountWith(table *godog.Table) error {
	cmd, err := assist.CreateInstance(new(queries.GetAccount), table)
	if err != nil {
		return err
	}

	f.account, f.err = f.app.GetAccount(context.Background(), *cmd.(*queries.GetAccount))

	return nil
}
