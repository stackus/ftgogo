package adapters

import (
	"context"

	"github.com/stackus/ftgogo/serviceapis/commonapi"
	deliverypb "github.com/stackus/ftgogo/serviceapis/deliveryapi/pb"
	"github.com/stackus/ftgogo/store-web/internal/domain"
)

type DeliveryGrpcRepository struct {
	client deliverypb.DeliveryServiceClient
}

var _ DeliveryRepository = (*DeliveryGrpcRepository)(nil)

func NewDeliveryGrpcRepository(client deliverypb.DeliveryServiceClient) *DeliveryGrpcRepository {
	return &DeliveryGrpcRepository{client: client}
}

func (r DeliveryGrpcRepository) FindDelivery(ctx context.Context, findDelivery FindDelivery) (*domain.Delivery, error) {
	resp, err := r.client.GetDelivery(ctx, &deliverypb.GetDeliveryRequest{DeliveryID: findDelivery.DeliveryID})
	if err != nil {
		return nil, err
	}

	return r.fromDeliveryProto(resp.Delivery), nil
}

func (r DeliveryGrpcRepository) FindCourier(ctx context.Context, findCourier FindCourier) (*domain.Courier, error) {
	resp, err := r.client.GetCourier(ctx, &deliverypb.GetCourierRequest{CourierID: findCourier.CourierID})
	if err != nil {
		return nil, err
	}

	return r.fromCourierProto(resp.Courier), nil
}

func (r DeliveryGrpcRepository) SetCourierAvailability(ctx context.Context, setCourierAvailability SetCourierAvailability) error {
	_, err := r.client.SetCourierAvailability(ctx, &deliverypb.SetCourierAvailabilityRequest{
		CourierID: setCourierAvailability.CourierID,
		Available: setCourierAvailability.Available,
	})
	return err
}

func (r DeliveryGrpcRepository) fromDeliveryProto(delivery *deliverypb.Delivery) *domain.Delivery {
	return &domain.Delivery{
		DeliveryID:        delivery.DeliveryID,
		RestaurantID:      delivery.RestaurantID,
		AssignedCourierID: delivery.AssignedCourierID,
		PickUpAddress:     commonapi.FromAddressProto(delivery.PickUpAddress),
		DeliveryAddress:   commonapi.FromAddressProto(delivery.DeliveryAddress),
		Status:            r.fromDeliveryStatusProto(delivery.Status),
		PickUpTime:        delivery.PickupTime.AsTime(),
		ReadyBy:           delivery.ReadyBy.AsTime(),
	}
}

func (r DeliveryGrpcRepository) fromCourierProto(courier *deliverypb.Courier) *domain.Courier {
	return &domain.Courier{
		CourierID: courier.CourierID,
		Plan:      r.fromPlanProto(courier.Plan),
		Available: courier.Available,
	}
}

func (r DeliveryGrpcRepository) fromPlanProto(plan *deliverypb.Plan) domain.Plan {
	p := make(domain.Plan, 0, len(plan.Actions))
	for _, action := range plan.Actions {
		p = append(p, domain.Action{
			DeliveryID: action.DeliveryID,
			ActionType: r.fromActionTypeProto(action.ActionType),
			Address:    commonapi.FromAddressProto(action.Address),
			When:       action.When.AsTime(),
		})
	}

	return p
}

func (r DeliveryGrpcRepository) fromDeliveryStatusProto(status deliverypb.DeliveryStatus) domain.DeliveryStatus {
	switch status {
	case deliverypb.DeliveryStatus_Scheduled:
		return domain.DeliveryScheduled
	case deliverypb.DeliveryStatus_Cancelled:
		return domain.DeliveryCancelled
	default:
		return domain.DeliveryPending
	}
}

func (r DeliveryGrpcRepository) fromActionTypeProto(actionType deliverypb.ActionType) domain.ActionType {
	switch actionType {
	case deliverypb.ActionType_PickUp:
		return domain.PickUp
	case deliverypb.ActionType_DropOff:
		return domain.DropOff
	default:
		return domain.PickUp
	}
}
