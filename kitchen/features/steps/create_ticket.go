package steps

import (
	"context"

	"github.com/cucumber/godog"
	"github.com/stackus/errors"

	"github.com/stackus/ftgogo/kitchen/internal/application/commands"
	"github.com/stackus/ftgogo/serviceapis/kitchenapi"
)

func (f *FeatureState) RegisterCreateTicketSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I (?:create|setup|have created) (?:a|the|another) ticket for order "([^"]*)" (?:and|at) restaurant "([^"]*)" with items$`, f.iCreateATicketForOrderAndRestaurant)
	ctx.Step(`^I (?:confirm|have confirmed) creating (?:a|the|another) ticket for order "([^"]*)"$`, f.iConfirmCreateATicketForOrder)
	ctx.Step(`^I (?:cancel|have cancell?ed) creating (?:a|the|another) ticket for order "([^"]*)"$`, f.iCancelCreateATicketForOrder)
}

func (f *FeatureState) iCreateATicketForOrderAndRestaurant(orderID, restaurantName string, table *godog.Table) error {
	restaurantID := f.restaurantIDs[restaurantName]

	menuItems, err := parseTableIntoMenuItems(table)
	if err != nil {
		return err
	}

	f.ticketID, err = f.app.CreateTicket(context.Background(), commands.CreateTicket{
		OrderID:      orderID,
		RestaurantID: restaurantID,
		LineItems:    menuItems,
	})

	f.ticketIDs[orderID] = f.ticketID

	return nil
}

func (f *FeatureState) iConfirmCreateATicketForOrder(orderID string) error {
	ticketID := f.ticketIDs[orderID]

	f.err = f.app.ConfirmCreateTicket(context.Background(), commands.ConfirmCreateTicket{TicketID: ticketID})

	return nil
}

func (f *FeatureState) iCancelCreateATicketForOrder(orderID string) error {
	ticketID := f.ticketIDs[orderID]

	f.err = f.app.CancelCreateTicket(context.Background(), commands.CancelCreateTicket{TicketID: ticketID})

	return nil
}

func parseTableIntoMenuItems(table *godog.Table) ([]kitchenapi.LineItem, error) {
	items, err := assist.CreateSlice(new(kitchenapi.LineItem), table)
	if err != nil {
		return nil, errors.Wrapf(errors.ErrUnknown, "error parsing menu items table: %w", err)
	}

	menuItems := make([]kitchenapi.LineItem, 0, len(items.([]*kitchenapi.LineItem)))

	for _, item := range items.([]*kitchenapi.LineItem) {
		menuItems = append(menuItems, *item)
	}

	return menuItems, nil
}
