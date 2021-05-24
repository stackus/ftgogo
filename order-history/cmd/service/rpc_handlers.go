package main

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/stackus/ftgogo/order-history/internal/application/queries"
	"github.com/stackus/ftgogo/serviceapis/orderapi"
	"github.com/stackus/ftgogo/serviceapis/orderhistoryapi/pb"
)

type rpcHandlers struct {
	app Application
	orderhistorypb.UnimplementedOrderHistoryServiceServer
}

var _ orderhistorypb.OrderHistoryServiceServer = (*rpcHandlers)(nil)

func newRpcHandlers(app Application) rpcHandlers {
	return rpcHandlers{app: app}
}

func (h rpcHandlers) GetConsumerOrderHistory(ctx context.Context, request *orderhistorypb.GetConsumerOrderHistoryRequest) (*orderhistorypb.GetConsumerOrderHistoryResponse, error) {
	var filters *queries.OrderHistoryFilters

	if request.Filter != nil {
		filters = &queries.OrderHistoryFilters{
			Since:    request.Filter.Since.AsTime(),
			Keywords: request.Filter.Keywords,
			Status:   orderapi.OrderState(int(request.Filter.Status)),
		}
	}

	results, err := h.app.Queries.GetConsumerOrderHistory.Handle(ctx, queries.GetConsumerOrderHistory{
		ConsumerID: request.ConsumerID,
		Filter:     filters,
		Next:       request.Next,
		Limit:      int(request.Limit),
	})
	if err != nil {
		return nil, err
	}

	orders := make([]*orderhistorypb.GetConsumerOrderHistoryResponseOrderHistory, len(results.Orders))
	for _, order := range results.Orders {
		orders = append(orders, &orderhistorypb.GetConsumerOrderHistoryResponseOrderHistory{
			OrderID:        order.OrderID,
			Status:         order.Status,
			RestaurantID:   order.RestaurantID,
			RestaurantName: order.RestaurantName,
			CreatedAt:      timestamppb.New(order.CreatedAt),
		})
	}

	return &orderhistorypb.GetConsumerOrderHistoryResponse{
		Orders: orders,
		Next:   results.Next,
	}, nil
}

func (h rpcHandlers) GetOrderHistory(ctx context.Context, request *orderhistorypb.GetOrderHistoryRequest) (*orderhistorypb.GetOrderHistoryResponse, error) {
	result, err := h.app.Queries.GetOrderHistory.Handle(ctx, queries.GetOrderHistory{OrderID: request.OrderID})
	if err != nil {
		return nil, err
	}

	return &orderhistorypb.GetOrderHistoryResponse{
		OrderID:        result.OrderID,
		Status:         result.Status,
		RestaurantID:   result.RestaurantID,
		RestaurantName: result.RestaurantName,
		CreatedAt:      timestamppb.New(result.CreatedAt),
	}, nil
}
