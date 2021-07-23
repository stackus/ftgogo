package steps

import (
	"context"
	"reflect"

	"github.com/cucumber/godog"
	"github.com/stackus/errors"

	"github.com/stackus/ftgogo/delivery/internal/application/queries"
	"github.com/stackus/ftgogo/delivery/internal/domain"
)

func (f *FeatureState) RegisterGetCourierSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I get the assigned courier for order "([^"]*)"$`, f.iGetTheAssignedCourierForOrder)

	ctx.Step(`^(?:ensure )?the returned courier will (pick-?up|drop-?off) the food at address$`, f.theReturnedCourierWillTheFoodAtAddress)
	ctx.Step(`^(?:ensure )?"([^"]*)" is (not(?: ))?the assigned courier$`, f.isTheAssignedCourier)
	ctx.Step(`^(?:ensure )?the returned courier is (not(?: ))?available$`, f.theReturnedCourierIsAvailable)
}

func (f *FeatureState) iGetTheAssignedCourierForOrder(orderID string) error {
	var err error

	f.delivery, err = f.app.GetDelivery(context.Background(), queries.GetDelivery{OrderID: orderID})
	if err != nil {
		return errors.Wrapf(errors.ErrUnknown, "error fetching order: %w", err)
	}

	if f.delivery != nil {
		f.assignedCourierID = f.delivery.AssignedCourierID
	}

	f.courier, f.err = f.app.GetCourier(context.Background(), queries.GetCourier{CourierID: f.assignedCourierID})

	return nil
}

func (f *FeatureState) theReturnedCourierWillTheFoodAtAddress(event string, table *godog.Table) error {
	actionType := domain.PickUp
	if event == "dropoff" || event == "drop-off" {
		actionType = domain.DropOff
	}
	expected, err := parseAddressFromTable(table)
	if err != nil {
		return err
	}

	if f.courier == nil {
		return errors.Wrap(errors.ErrNotFound, "expected courier to not be nil")
	}

	var action domain.Action
	for _, a := range f.courier.Plan {
		if a.ActionType == actionType {
			action = a
			break
		}
	}

	if action.DeliveryID == "" {
		return errors.Wrapf(errors.ErrNotFound, "courier is missing a %s action", event)
	}

	if !reflect.DeepEqual(expected, action.Address) {
		return errors.Wrapf(errors.ErrInvalidArgument, "does not match expected: %v: got: %v", expected, action.Address)
	}

	return nil
}

func (f *FeatureState) isTheAssignedCourier(courierName, neg string) error {
	courierID := f.courierIDs[courierName]

	if neg == "not" && f.assignedCourierID == courierID {
		return errors.Wrapf(errors.ErrInvalidArgument, "expected assigned courier to not be %s", courierName)
	} else if neg == "" && f.assignedCourierID != courierID {
		return errors.Wrapf(errors.ErrInvalidArgument, "expected assigned courier to be %s", courierName)
	}

	return nil
}

func (f *FeatureState) theReturnedCourierIsAvailable(neg string) error {
	if f.courier == nil {
		return errors.Wrap(errors.ErrNotFound, "Expected courier to not be nil")
	}

	if neg == "not" && f.courier.Available {
		return errors.Wrap(errors.ErrInvalidArgument, "expected courier to not be available")
	} else if neg == "" && !f.courier.Available {
		return errors.Wrap(errors.ErrInvalidArgument, "expected courier to be available")
	}

	return nil
}
