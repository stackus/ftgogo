package adapters

import (
	"context"

	"github.com/stackus/ftgogo/serviceapis/orderapi"
	"github.com/stackus/ftgogo/serviceapis/orderapi/pb"
	"github.com/stackus/ftgogo/store-web/internal/domain"
)

type OrderGRPCRepository struct {
	client orderpb.OrderServiceClient
}

var _ OrderRepository = (*OrderGRPCRepository)(nil)

func NewOrderGrpcClient(client orderpb.OrderServiceClient) *OrderGRPCRepository {
	return &OrderGRPCRepository{client: client}
}

func (r OrderGRPCRepository) Find(ctx context.Context, findOrder FindOrder) (*domain.Order, error) {
	resp, err := r.client.GetOrder(ctx, &orderpb.GetOrderRequest{OrderID: findOrder.OrderID})
	if err != nil {
		return nil, err
	}

	return r.fromOrderProto(resp.Order), nil
}

func (r OrderGRPCRepository) Cancel(ctx context.Context, cancelOrder CancelOrder) (orderapi.OrderState, error) {
	resp, err := r.client.CancelOrder(ctx, &orderpb.CancelOrderRequest{
		OrderID: cancelOrder.OrderID,
	})
	if err != nil {
		return orderapi.UnknownOrderState, err
	}
	return orderapi.FromOrderStateProto(resp.Status), err
}

func (r OrderGRPCRepository) fromOrderProto(order *orderpb.Order) *domain.Order {
	return &domain.Order{
		OrderID: order.OrderID,
		Total:   int(order.OrderTotal),
		Status:  orderapi.FromOrderStateProto(order.Status),
	}
}
