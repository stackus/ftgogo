package steps

import (
	"context"
	"time"

	"github.com/cucumber/godog"

	"github.com/stackus/ftgogo/kitchen/internal/application/commands"
)

func (f *FeatureState) RegisterAcceptTicketSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I (?:accept|have accepted) (?:a|the|another) ticket for order "([^"]*)" will be ready in (\d+) minutes$`, f.iAcceptThatTicketWillBeReadyInMinutesForOrder)
}

func (f *FeatureState) iAcceptThatTicketWillBeReadyInMinutesForOrder(orderID string, minutes int) error {
	ticketID := f.ticketIDs[orderID]

	f.err = f.app.AcceptTicket(context.Background(), commands.AcceptTicket{
		TicketID: ticketID,
		ReadyBy:  time.Now().Add(time.Minute * time.Duration(minutes)),
	})

	return nil
}
