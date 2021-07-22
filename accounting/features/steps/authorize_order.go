package steps

import (
	"context"

	"github.com/cucumber/godog"
	"github.com/google/uuid"

	"github.com/stackus/ftgogo/accounting/internal/application/commands"
)

func (f *FeatureState) RegisterAuthorizeOrderSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I authorize an order totaling \$(\d+)\.(\d+) for "([^"]*)"$`, f.iAuthorizeAnOrderTotalingFor)
	ctx.Step(`^I expect it be authorized$`, f.iExpectTheCommandToSucceed)
	ctx.Step(`^I don\'t expect it to be authorized$`, f.iExpectTheCommandToFail)
}

func (f *FeatureState) iAuthorizeAnOrderTotalingFor(dollars, cents int, consumerName string) error {
	orderID := uuid.New().String()
	consumerID := f.accountNames[consumerName]

	cmd := commands.AuthorizeOrder{
		ConsumerID: consumerID,
		OrderID:    orderID,
		OrderTotal: dollars*10 + cents,
	}

	f.err = f.app.AuthorizeOrder(context.Background(), cmd)

	return nil
}
