package steps

import (
	"context"

	"github.com/cucumber/godog"

	"github.com/stackus/ftgogo/kitchen/internal/application/commands"
)

func (f *FeatureState) RegisterCancelTicketSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I (?:begin|began|have begun) (?:to cancel|cancelling) (?:a|the|another) ticket for order "([^"]*)"$`, f.iBeginToCancelATicketForOrder)
	ctx.Step(`^I (?:undo|have undone) cancell?ing (?:a|the|another) ticket for order "([^"]*)"$`, f.iUndoCancelATicketForOrder)
	ctx.Step(`^I (?:confirm|have confirmed) cancell?ing (?:a|the|another) ticket for order "([^"]*)"$`, f.iConfirmCancelATicketForOrder)
}

func (f *FeatureState) iBeginToCancelATicketForOrder(orderID string) error {
	ticketID := f.ticketIDs[orderID]

	f.err = f.app.BeginCancelTicket(context.Background(), commands.BeginCancelTicket{
		TicketID: ticketID,
	})

	return nil
}

func (f *FeatureState) iConfirmCancelATicketForOrder(orderID string) error {
	ticketID := f.ticketIDs[orderID]

	f.err = f.app.ConfirmCancelTicket(context.Background(), commands.ConfirmCancelTicket{
		TicketID: ticketID,
	})

	return nil
}

func (f *FeatureState) iUndoCancelATicketForOrder(orderID string) error {
	ticketID := f.ticketIDs[orderID]

	f.err = f.app.UndoCancelTicket(context.Background(), commands.UndoCancelTicket{
		TicketID: ticketID,
	})

	return nil
}
