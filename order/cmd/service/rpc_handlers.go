package main

import (
	"context"

	"github.com/stackus/ftgogo/order/internal/application/commands"
	"github.com/stackus/ftgogo/order/internal/application/queries"
	"github.com/stackus/ftgogo/order/internal/domain"
	"github.com/stackus/ftgogo/serviceapis/commonapi"
	"github.com/stackus/ftgogo/serviceapis/commonapi/pb"
	"github.com/stackus/ftgogo/serviceapis/orderapi"
	"github.com/stackus/ftgogo/serviceapis/orderapi/pb"
)

type rpcHandlers struct {
	app Application
	orderpb.UnimplementedOrderServiceServer
}

var _ orderpb.OrderServiceServer = (*rpcHandlers)(nil)

func newRpcHandlers(app Application) rpcHandlers {
	return rpcHandlers{app: app}
}

func (h rpcHandlers) CreateOrder(ctx context.Context, request *orderpb.CreateOrderRequest) (*orderpb.CreateOrderResponse, error) {
	lineItems := make(map[string]int, len(request.LineItems.Items))
	for s, i := range request.LineItems.Items {
		lineItems[s] = int(i)
	}

	orderID, err := h.app.Commands.CreateOrder.Handle(ctx, commands.CreateOrder{
		ConsumerID:   request.ConsumerID,
		RestaurantID: request.RestaurantID,
		DeliverAt:    request.DeliverAt.AsTime(),
		DeliverTo: commonapi.Address{
			Street1: request.DeliverTo.Street1,
			Street2: request.DeliverTo.Street2,
			City:    request.DeliverTo.City,
			State:   request.DeliverTo.State,
			Zip:     request.DeliverTo.Zip,
		},
		LineItems: lineItems,
	})
	if err != nil {
		return nil, err
	}

	return &orderpb.CreateOrderResponse{OrderID: orderID}, nil
}

func (h rpcHandlers) GetOrder(ctx context.Context, request *orderpb.GetOrderRequest) (*orderpb.GetOrderResponse, error) {
	order, err := h.app.Queries.GetOrder.Handle(ctx, queries.GetOrder{OrderID: request.OrderID})
	if err != nil {
		return nil, err
	}

	return &orderpb.GetOrderResponse{
		Order: h.toOrderProto(order),
	}, nil
}

func (h rpcHandlers) CancelOrder(ctx context.Context, request *orderpb.CancelOrderRequest) (*orderpb.CancelOrderResponse, error) {
	status, err := h.app.Commands.StartCancelOrderSaga.Handle(ctx, commands.StartCancelOrderSaga{OrderID: request.OrderID})
	if err != nil {
		return nil, err
	}

	return &orderpb.CancelOrderResponse{Status: h.toOrderStateProto(status)}, nil
}

func (h rpcHandlers) ReviseOrder(ctx context.Context, request *orderpb.ReviseOrderRequest) (*orderpb.ReviseOrderResponse, error) {
	status, err := h.app.Commands.StartReviseOrderSaga.Handle(ctx, commands.StartReviseOrderSaga{
		OrderID:           request.OrderID,
		RevisedQuantities: h.fromMenuItemQuantitiesProto(request.RevisedQuantities),
	})
	if err != nil {
		return nil, err
	}

	return &orderpb.ReviseOrderResponse{Status: h.toOrderStateProto(status)}, nil
}

func (h rpcHandlers) toOrderProto(order *domain.Order) *orderpb.Order {
	return &orderpb.Order{
		OrderID:      order.ID(),
		RestaurantID: order.RestaurantID,
		ConsumerID:   order.ConsumerID,
		OrderTotal:   int64(order.OrderTotal()),
		Status:       h.toOrderStateProto(order.State),
	}
}

func (h rpcHandlers) toOrderStateProto(state orderapi.OrderState) orderpb.OrderState {
	switch state {
	case orderapi.ApprovalPending:
		return orderpb.OrderState_ApprovalPending
	case orderapi.Approved:
		return orderpb.OrderState_Approved
	case orderapi.CancelPending:
		return orderpb.OrderState_CancelPending
	case orderapi.Cancelled:
		return orderpb.OrderState_Cancelled
	case orderapi.RevisionPending:
		return orderpb.OrderState_RevisionPending
	case orderapi.Rejected:
		return orderpb.OrderState_Rejected
	default:
		return orderpb.OrderState_Unknown
	}
}

func (h rpcHandlers) toMenuItemsQuantitiesProto(quantities commonapi.MenuItemQuantities) *commonpb.MenuItemQuantities {
	lineItems := make(map[string]int64, len(quantities))
	for itemID, qty := range quantities {
		lineItems[itemID] = int64(qty)
	}

	return &commonpb.MenuItemQuantities{Items: lineItems}
}

func (h rpcHandlers) fromMenuItemQuantitiesProto(quantities *commonpb.MenuItemQuantities) commonapi.MenuItemQuantities {
	lineItems := make(commonapi.MenuItemQuantities, len(quantities.Items))
	for itemID, qty := range quantities.Items {
		lineItems[itemID] = int(qty)
	}

	return lineItems
}