package handlers

import (
	"context"

	"google.golang.org/grpc"

	"github.com/stackus/ftgogo/delivery/internal/application"
	"github.com/stackus/ftgogo/delivery/internal/application/commands"
	"github.com/stackus/ftgogo/delivery/internal/application/queries"
	"github.com/stackus/ftgogo/serviceapis/deliveryapi/pb"
)

type RpcHandlers struct {
	app application.Service
	deliverypb.UnimplementedDeliveryServiceServer
}

var _ deliverypb.DeliveryServiceServer = (*RpcHandlers)(nil)

func NewRpcHandlers(app application.Service) RpcHandlers {
	return RpcHandlers{app: app}
}

func (h RpcHandlers) Mount(registrar grpc.ServiceRegistrar) {
	deliverypb.RegisterDeliveryServiceServer(registrar, h)
}

func (h RpcHandlers) SetCourierAvailability(ctx context.Context, request *deliverypb.SetCourierAvailabilityRequest) (*deliverypb.SetCourierAvailabilityResponse, error) {
	err := h.app.Commands.SetCourierAvailability.Handle(ctx, commands.SetCourierAvailability{
		CourierID: request.CourierID,
		Available: request.Available,
	})
	if err != nil {
		return nil, err
	}

	return &deliverypb.SetCourierAvailabilityResponse{Available: request.Available}, nil
}

func (h RpcHandlers) GetDeliveryStatus(ctx context.Context, request *deliverypb.GetDeliveryStatusRequest) (*deliverypb.GetDeliveryStatusResponse, error) {
	status, err := h.app.Queries.GetDeliveryStatus.Handle(ctx, queries.GetDeliveryStatus{DeliveryID: request.DeliveryID})
	if err != nil {
		return nil, err
	}

	return &deliverypb.GetDeliveryStatusResponse{
		DeliveryID:        status.ID,
		AssignedCourierID: status.AssignedCourier,
		CourierActions:    status.CourierActions,
		Status:            status.Status,
	}, nil
}
