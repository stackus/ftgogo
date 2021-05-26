package adapters

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/stackus/ftgogo/serviceapis/commonapi"
	"github.com/stackus/ftgogo/serviceapis/orderapi"
	"github.com/stackus/ftgogo/serviceapis/orderapi/pb"
	"github.com/stackus/ftgogo/web-bff/internal/domain"
)

type OrderGRPCRepository struct {
	client orderpb.OrderServiceClient
}

var _ domain.OrderRepository = (*OrderGRPCRepository)(nil)

func NewOrderGrpcClient(client orderpb.OrderServiceClient) *OrderGRPCRepository {
	return &OrderGRPCRepository{client: client}
}

func (r OrderGRPCRepository) Create(ctx context.Context, createOrder domain.CreateOrder) (string, error) {
	res, err := r.client.CreateOrder(
		ctx, &orderpb.CreateOrderRequest{
			ConsumerID:   createOrder.ConsumerID,
			RestaurantID: createOrder.RestaurantID,
			DeliverTo:    commonapi.ToAddressProto(createOrder.DeliverTo),
			DeliverAt:    timestamppb.New(createOrder.DeliverAt),
			LineItems:    commonapi.ToMenuItemQuantitiesProto(createOrder.LineItems),
		},
	)
	if err != nil {
		return "", err
	}

	return res.OrderID, nil
}

func (r OrderGRPCRepository) Find(ctx context.Context, findOrder domain.FindOrder) (*domain.Order, error) {
	resp, err := r.client.GetOrder(ctx, &orderpb.GetOrderRequest{OrderID: findOrder.OrderID})
	if err != nil {
		return nil, err
	}

	return r.fromOrderProto(resp.Order), nil
}

func (r OrderGRPCRepository) Revise(ctx context.Context, reviseOrder domain.ReviseOrder) (orderapi.OrderState, error) {
	resp, err := r.client.ReviseOrder(ctx, &orderpb.ReviseOrderRequest{
		OrderID:           reviseOrder.OrderID,
		RevisedQuantities: commonapi.ToMenuItemQuantitiesProto(reviseOrder.RevisedQuantities),
	})
	if err != nil {
		return orderapi.UnknownOrderState, err
	}
	return orderapi.FromOrderStateProto(resp.Status), err
}

func (r OrderGRPCRepository) Cancel(ctx context.Context, cancelOrder domain.CancelOrder) (orderapi.OrderState, error) {
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
