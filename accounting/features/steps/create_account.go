package steps

import (
	"context"

	"github.com/cucumber/godog"
	"github.com/google/uuid"

	"github.com/stackus/ftgogo/accounting/internal/application/commands"
)

func (f *FeatureState) RegisterCreateAccountSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I create an account for the consumer "([^"]*)"$`, f.iCreateAnAccountForTheConsumer)
}

func (f *FeatureState) iCreateAnAccountForTheConsumer(consumerName string) error {
	var consumerID string

	consumerID = f.accountNames[consumerName]
	if consumerID == "" {
		consumerID = uuid.New().String()
		f.accountNames[consumerName] = consumerID
	}

	cmd := commands.CreateAccount{
		ConsumerID: consumerID,
		Name:       consumerName,
	}

	f.err = f.app.CreateAccount(context.Background(), cmd)

	return nil
}
