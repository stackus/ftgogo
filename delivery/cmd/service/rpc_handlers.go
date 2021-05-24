package main

import (
	"context"

	"github.com/stackus/ftgogo/delivery/internal/application/commands"
	"github.com/stackus/ftgogo/delivery/internal/application/queries"
	"github.com/stackus/ftgogo/serviceapis/deliveryapi/pb"
)

type rpcHandlers struct {
	app Application
	deliverypb.UnimplementedDeliveryServiceServer
}

var _ deliverypb.DeliveryServiceServer = (*rpcHandlers)(nil)

func newRpcHandlers(app Application) rpcHandlers {
	return rpcHandlers{app: app}
}

func (h rpcHandlers) SetCourierAvailability(ctx context.Context, request *deliverypb.SetCourierAvailabilityRequest) (*deliverypb.SetCourierAvailabilityResponse, error) {
	err := h.app.Commands.SetCourierAvailability.Handle(ctx, commands.SetCourierAvailability{
		CourierID: request.CourierID,
		Available: request.Available,
	})
	if err != nil {
		return nil, err
	}

	return &deliverypb.SetCourierAvailabilityResponse{Available: request.Available}, nil
}

func (h rpcHandlers) GetDeliveryStatus(ctx context.Context, request *deliverypb.GetDeliveryStatusRequest) (*deliverypb.GetDeliveryStatusResponse, error) {
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
