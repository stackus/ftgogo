package adapters

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/stackus/ftgogo/customer-web/internal/domain"
	"github.com/stackus/ftgogo/serviceapis/commonapi"
	"github.com/stackus/ftgogo/serviceapis/orderapi"
	"github.com/stackus/ftgogo/serviceapis/orderapi/pb"
)

type OrderGRPCRepository struct {
	client orderpb.OrderServiceClient
}

var _ OrderRepository = (*OrderGRPCRepository)(nil)

func NewOrderGrpcClient(client orderpb.OrderServiceClient) *OrderGRPCRepository {
	return &OrderGRPCRepository{client: client}
}

func (r OrderGRPCRepository) Create(ctx context.Context, createOrder CreateOrder) (string, error) {
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

func (r OrderGRPCRepository) Find(ctx context.Context, findOrder FindOrder) (*domain.Order, error) {
	resp, err := r.client.GetOrder(ctx, &orderpb.GetOrderRequest{OrderID: findOrder.OrderID})
	if err != nil {
		return nil, err
	}

	return r.fromOrderProto(resp.Order), nil
}

func (r OrderGRPCRepository) Revise(ctx context.Context, reviseOrder ReviseOrder) (orderapi.OrderState, error) {
	resp, err := r.client.ReviseOrder(ctx, &orderpb.ReviseOrderRequest{
		OrderID:           reviseOrder.OrderID,
		RevisedQuantities: commonapi.ToMenuItemQuantitiesProto(reviseOrder.RevisedQuantities),
	})
	if err != nil {
		return orderapi.UnknownOrderState, err
	}
	return orderapi.FromOrderStateProto(resp.Status), err
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
