package steps

import (
	"github.com/cucumber/godog"
	"github.com/rdumont/assistdog"
	_ "github.com/stackus/edat-msgpack"
	"github.com/stackus/errors"

	"github.com/stackus/ftgogo/delivery/internal/adapters"
	"github.com/stackus/ftgogo/delivery/internal/application"
	"github.com/stackus/ftgogo/delivery/internal/domain"
	"github.com/stackus/ftgogo/serviceapis"
	"github.com/stackus/ftgogo/serviceapis/commonapi"
)

type ConsumerJson struct {
	ID        string
	Name      string
	Addresses map[string]*commonapi.Address
}

type FeatureState struct {
	app               application.ServiceApplication
	courier           *domain.Courier
	delivery          *domain.Delivery
	assignedCourierID string
	courierIDs        map[string]string
	restaurantIDs     map[string]string
	err               error
}

func NewFeatureState() *FeatureState {
	f := &FeatureState{}
	f.Reset()

	return f
}

var assist = assistdog.NewDefault()

func init() {
	serviceapis.RegisterTypes()
}

func (f *FeatureState) Reset() {
	f.courier = nil
	f.delivery = nil
	f.assignedCourierID = ""
	f.courierIDs = make(map[string]string)
	f.restaurantIDs = make(map[string]string)
	f.err = nil

	courierRepo := adapters.NewCourierInmemRepository()
	deliveryRepo := adapters.NewDeliveryInmemRepository()
	restaurantRepo := adapters.NewRestaurantInmemRepository()
	f.app = application.NewServiceApplication(courierRepo, deliveryRepo, restaurantRepo)
}

func (f *FeatureState) RegisterCommonSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I expect the (?:request|command|query) to fail$`, f.iExpectTheCommandToFail)
	ctx.Step(`^I expect the (?:request|command|query) to succeed$`, f.iExpectTheCommandToSucceed)

	ctx.Step(`^(?:ensure |expect )?the returned error message is "([^"]*)"$`, f.theReturnedErrorMessageIs)
}

func (f *FeatureState) iExpectTheCommandToFail() error {
	if f.err == nil {
		return errors.Wrap(errors.ErrUnknown, "expected error to not be nil")
	}

	return nil
}

func (f *FeatureState) iExpectTheCommandToSucceed() error {
	if f.err != nil {
		return errors.Wrapf(f.err, "expected error to be nil: got %w", f.err)
	}

	return nil
}

func (f *FeatureState) theReturnedErrorMessageIs(errorMsg string) error {
	if f.err == nil {
		return errors.Wrap(errors.ErrUnknown, "expected error to not be nil")
	}

	if errorMsg != f.err.Error() {
		return errors.Wrapf(errors.ErrInvalidArgument, "expected: %s: got: %s", errorMsg, f.err.Error())
	}

	return nil
}

func parseAddressFromTable(table *godog.Table) (*commonapi.Address, error) {
	address, err := assist.CreateInstance(new(commonapi.Address), table)
	if err != nil {
		return nil, errors.Wrapf(errors.ErrUnknown, "error parsing address table: %w", err)
	}

	return address.(*commonapi.Address), nil
}
