package handlers

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/stackus/ftgogo/delivery/internal/application"
	"github.com/stackus/ftgogo/delivery/internal/application/commands"
	"github.com/stackus/ftgogo/delivery/internal/application/queries"
	"github.com/stackus/ftgogo/delivery/internal/domain"
	"github.com/stackus/ftgogo/serviceapis/commonapi"
	"github.com/stackus/ftgogo/serviceapis/deliveryapi/pb"
)

type RpcHandlers struct {
	app application.ServiceApplication
	deliverypb.UnimplementedDeliveryServiceServer
}

var _ deliverypb.DeliveryServiceServer = (*RpcHandlers)(nil)

func NewRpcHandlers(app application.ServiceApplication) RpcHandlers {
	return RpcHandlers{app: app}
}

func (h RpcHandlers) Mount(registrar grpc.ServiceRegistrar) {
	deliverypb.RegisterDeliveryServiceServer(registrar, h)
}

func (h RpcHandlers) GetCourier(ctx context.Context, request *deliverypb.GetCourierRequest) (*deliverypb.GetCourierResponse, error) {
	courier, err := h.app.GetCourier(ctx, queries.GetCourier{CourierID: request.CourierID})
	if err != nil {
		return nil, err
	}
	return &deliverypb.GetCourierResponse{
		Courier: h.toCourierProto(courier),
	}, nil
}

func (h RpcHandlers) SetCourierAvailability(ctx context.Context, request *deliverypb.SetCourierAvailabilityRequest) (*deliverypb.SetCourierAvailabilityResponse, error) {
	err := h.app.SetCourierAvailability(ctx, commands.SetCourierAvailability{
		CourierID: request.CourierID,
		Available: request.Available,
	})
	if err != nil {
		return nil, err
	}

	return &deliverypb.SetCourierAvailabilityResponse{Available: request.Available}, nil
}

func (h RpcHandlers) GetDelivery(ctx context.Context, request *deliverypb.GetDeliveryRequest) (*deliverypb.GetDeliveryResponse, error) {
	delivery, err := h.app.GetDelivery(ctx, queries.GetDelivery{OrderID: request.DeliveryID})
	if err != nil {
		return nil, err
	}
	return &deliverypb.GetDeliveryResponse{
		Delivery: h.toDeliveryProto(delivery),
	}, nil
}

func (h RpcHandlers) toDeliveryProto(delivery *domain.Delivery) *deliverypb.Delivery {
	return &deliverypb.Delivery{
		DeliveryID:        delivery.DeliveryID,
		RestaurantID:      delivery.RestaurantID,
		AssignedCourierID: delivery.AssignedCourierID,
		Status:            h.toDeliveryStatus(delivery.Status),
		PickUpAddress:     commonapi.ToAddressProto(delivery.PickUpAddress),
		DeliveryAddress:   commonapi.ToAddressProto(delivery.DeliveryAddress),
		PickupTime:        timestamppb.New(delivery.PickUpTime),
		ReadyBy:           timestamppb.New(delivery.ReadyBy),
	}
}

func (h RpcHandlers) toCourierProto(courier *domain.Courier) *deliverypb.Courier {
	return &deliverypb.Courier{
		CourierID: courier.CourierID,
		Plan:      h.toPlanProto(courier.Plan),
		Available: courier.Available,
	}
}

func (h RpcHandlers) toPlanProto(plan domain.Plan) *deliverypb.Plan {
	actions := make([]*deliverypb.Action, 0, len(plan))

	for _, action := range plan {
		actions = append(actions, &deliverypb.Action{
			DeliveryID: action.DeliveryID,
			ActionType: h.toActionTypeProto(action.ActionType),
			Address:    commonapi.ToAddressProto(action.Address),
			When:       timestamppb.New(action.When),
		})
	}

	return &deliverypb.Plan{
		Actions: actions,
	}
}

func (h RpcHandlers) toDeliveryStatus(status domain.DeliveryStatus) deliverypb.DeliveryStatus {
	switch status {
	case domain.DeliveryScheduled:
		return deliverypb.DeliveryStatus_Scheduled
	case domain.DeliveryCancelled:
		return deliverypb.DeliveryStatus_Cancelled
	default:
		return deliverypb.DeliveryStatus_Pending
	}
}

func (h RpcHandlers) toActionTypeProto(actionType domain.ActionType) deliverypb.ActionType {
	switch actionType {
	case domain.PickUp:
		return deliverypb.ActionType_PickUp
	case domain.DropOff:
		return deliverypb.ActionType_DropOff
	default:
		return deliverypb.ActionType_PickUp
	}
}
