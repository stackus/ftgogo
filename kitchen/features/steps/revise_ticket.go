package steps

import (
	"context"

	"github.com/cucumber/godog"

	"github.com/stackus/ftgogo/kitchen/internal/application/commands"
)

func (f *FeatureState) RegisterReviseTicketSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I (?:begin|began|have begun) (?:to revise|revising) (?:a|the|another) ticket for order "([^"]*)"$`, f.iBeginToReviseATicketForOrder)
	ctx.Step(`^I (?:undo|have undone) revising (?:a|the|another) ticket for order "([^"]*)"$`, f.iUndoReviseATicketForOrder)
	ctx.Step(`^I (?:confirm|have confirmed) revising (?:a|the|another) ticket for order "([^"]*)"$`, f.iConfirmReviseATicketForOrder)
}

func (f *FeatureState) iBeginToReviseATicketForOrder(orderID string) error {
	ticketID := f.ticketIDs[orderID]

	f.err = f.app.BeginReviseTicket(context.Background(), commands.BeginReviseTicket{
		TicketID: ticketID,
	})

	return nil
}

func (f *FeatureState) iConfirmReviseATicketForOrder(orderID string) error {
	ticketID := f.ticketIDs[orderID]

	f.err = f.app.ConfirmReviseTicket(context.Background(), commands.ConfirmReviseTicket{
		TicketID: ticketID,
	})

	return nil
}

func (f *FeatureState) iUndoReviseATicketForOrder(orderID string) error {
	ticketID := f.ticketIDs[orderID]

	f.err = f.app.UndoReviseTicket(context.Background(), commands.UndoReviseTicket{
		TicketID: ticketID,
	})

	return nil
}
