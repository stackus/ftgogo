package steps

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/cucumber/godog"
	"github.com/stackus/errors"

	"github.com/stackus/ftgogo/kitchen/internal/application/queries"
	"github.com/stackus/ftgogo/kitchen/internal/domain"
)

func (f *FeatureState) RegisterGetTicketSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I (?:get|fetch|request) (?:a|the|another) ticket for order "([^"]*)"$`, f.iGetATicketForOrder)

	ctx.Step(`^(?:ensure |expect )?the returned ticket matches:$`, f.theReturnedTicketMatches)
	ctx.Step(`^(?:ensure |expect )?the returned ticket status is "([^"]*)"$`, f.theReturnedTicketStatusIs)
}

func (f *FeatureState) iGetATicketForOrder(orderID string) error {
	ticketID := f.ticketIDs[orderID]

	f.ticket, f.err = f.app.GetTicket(context.Background(), queries.GetTicket{TicketID: ticketID})

	return nil
}

func (f *FeatureState) theReturnedTicketMatches(doc *godog.DocString) error {
	var err error
	var expected *domain.Ticket

	if err = json.Unmarshal([]byte(doc.Content), &expected); err != nil {
		return errors.Wrap(errors.ErrUnprocessableEntity, err.Error())
	}

	if f.ticket == nil {
		return errors.Wrap(errors.ErrNotFound, "Expected courier to not be nil")
	}

	if !reflect.DeepEqual(expected, f.ticket) {
		return errors.Wrapf(errors.ErrInvalidArgument, "does not match expected: %v: got: %v", expected, f.ticket)
	}

	return nil
}

func (f *FeatureState) theReturnedTicketStatusIs(expected string) error {
	if f.ticket == nil {
		return errors.Wrap(errors.ErrNotFound, "Expected courier to not be nil")
	}

	if f.ticket.State.String() != expected {
		return errors.Wrapf(errors.ErrInvalidArgument, "does not match expected: %v: got: %v", expected, f.ticket.State.String())
	}

	return nil
}
