package handlers

import (
	"context"

	"google.golang.org/grpc"

	"github.com/stackus/ftgogo/order/internal/application"
	"github.com/stackus/ftgogo/order/internal/application/commands"
	"github.com/stackus/ftgogo/order/internal/application/queries"
	"github.com/stackus/ftgogo/order/internal/domain"
	"github.com/stackus/ftgogo/serviceapis/commonapi"
	"github.com/stackus/ftgogo/serviceapis/orderapi"
	"github.com/stackus/ftgogo/serviceapis/orderapi/pb"
)

type RpcHandlers struct {
	app application.ServiceApplication
	orderpb.UnimplementedOrderServiceServer
}

var _ orderpb.OrderServiceServer = (*RpcHandlers)(nil)

func NewRpcHandlers(app application.ServiceApplication) RpcHandlers {
	return RpcHandlers{app: app}
}

func (h RpcHandlers) Mount(registrar grpc.ServiceRegistrar) {
	orderpb.RegisterOrderServiceServer(registrar, h)
}

func (h RpcHandlers) CreateOrder(ctx context.Context, request *orderpb.CreateOrderRequest) (*orderpb.CreateOrderResponse, error) {
	lineItems := make(map[string]int, len(request.LineItems.Items))
	for s, i := range request.LineItems.Items {
		lineItems[s] = int(i)
	}

	orderID, err := h.app.CreateOrder(ctx, commands.CreateOrder{
		ConsumerID:   request.ConsumerID,
		RestaurantID: request.RestaurantID,
		DeliverAt:    request.DeliverAt.AsTime(),
		DeliverTo:    commonapi.FromAddressProto(request.DeliverTo),
		LineItems:    lineItems,
	})
	if err != nil {
		return nil, err
	}

	return &orderpb.CreateOrderResponse{OrderID: orderID}, nil
}

func (h RpcHandlers) GetOrder(ctx context.Context, request *orderpb.GetOrderRequest) (*orderpb.GetOrderResponse, error) {
	order, err := h.app.GetOrder(ctx, queries.GetOrder{OrderID: request.OrderID})
	if err != nil {
		return nil, err
	}

	return &orderpb.GetOrderResponse{
		Order: h.toOrderProto(order),
	}, nil
}

func (h RpcHandlers) CancelOrder(ctx context.Context, request *orderpb.CancelOrderRequest) (*orderpb.CancelOrderResponse, error) {
	status, err := h.app.StartCancelOrderSaga(ctx, commands.StartCancelOrderSaga{OrderID: request.OrderID})
	if err != nil {
		return nil, err
	}

	return &orderpb.CancelOrderResponse{Status: orderapi.ToOrderStateProto(status)}, nil
}

func (h RpcHandlers) ReviseOrder(ctx context.Context, request *orderpb.ReviseOrderRequest) (*orderpb.ReviseOrderResponse, error) {
	status, err := h.app.StartReviseOrderSaga(ctx, commands.StartReviseOrderSaga{
		OrderID:           request.OrderID,
		RevisedQuantities: commonapi.FromMenuItemQuantitiesProto(request.RevisedQuantities),
	})
	if err != nil {
		return nil, err
	}

	return &orderpb.ReviseOrderResponse{Status: orderapi.ToOrderStateProto(status)}, nil
}

func (h RpcHandlers) toOrderProto(order *domain.Order) *orderpb.Order {
	return &orderpb.Order{
		OrderID:      order.ID(),
		RestaurantID: order.RestaurantID,
		ConsumerID:   order.ConsumerID,
		OrderTotal:   int64(order.OrderTotal()),
		Status:       orderapi.ToOrderStateProto(order.State),
	}
}
