package adapters

import (
	"context"

	"github.com/stackus/errors"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/stackus/ftgogo/serviceapis/commonapi"
	"github.com/stackus/ftgogo/serviceapis/commonapi/pb"
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
			DeliverTo: &commonpb.Address{
				Street1: createOrder.DeliverTo.Street1,
				Street2: createOrder.DeliverTo.Street2,
				City:    createOrder.DeliverTo.City,
				State:   createOrder.DeliverTo.State,
				Zip:     createOrder.DeliverTo.Zip,
			},
			DeliverAt: timestamppb.New(createOrder.DeliverAt),
			LineItems: r.toMenuItemQuantitiesProto(createOrder.LineItems),
		},
	)
	if err != nil {
		return "", err
	}

	return res.OrderID, nil
}

func (r OrderGRPCRepository) Find(ctx context.Context, orderID string) (*domain.Order, error) {
	resp, err := r.client.GetOrder(ctx, &orderpb.GetOrderRequest{OrderID: orderID})
	if err != nil {
		return nil, err
	}

	return r.fromOrderProto(resp.Order), nil
}

func (r OrderGRPCRepository) Revise(ctx context.Context, reviseOrder domain.ReviseOrder) error {
	var err error

	var resp *orderpb.GetOrderResponse
	resp, err = r.client.GetOrder(ctx, &orderpb.GetOrderRequest{OrderID: reviseOrder.OrderID})
	if err != nil {
		return err
	}

	if resp.Order.ConsumerID != reviseOrder.ConsumerID {
		return errors.Wrap(errors.ErrPermissionDenied, "you are not permitted to revise orders other than your own")
	}

	_, err = r.client.ReviseOrder(ctx, &orderpb.ReviseOrderRequest{
		OrderID:           reviseOrder.OrderID,
		RevisedQuantities: r.toMenuItemQuantitiesProto(reviseOrder.RevisedQuantities),
	})

	return err
}

func (r OrderGRPCRepository) Cancel(ctx context.Context, orderID string) error {
	panic("implement me")
}

func (r OrderGRPCRepository) toMenuItemQuantitiesProto(quantities commonapi.MenuItemQuantities) *commonpb.MenuItemQuantities {
	lineItems := make(map[string]int64, len(quantities))
	for itemID, qty := range quantities {
		lineItems[itemID] = int64(qty)
	}

	return &commonpb.MenuItemQuantities{Items: lineItems}
}

func (r OrderGRPCRepository) fromOrderProto(order *orderpb.Order) *domain.Order {
	return &domain.Order{
		OrderID: order.OrderID,
		Total:   int(order.OrderTotal),
		Status:  order.Status.String(),
	}
}
