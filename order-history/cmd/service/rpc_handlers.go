package main

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/stackus/ftgogo/order-history/internal/application/queries"
	"github.com/stackus/ftgogo/order-history/internal/domain"
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

func (h rpcHandlers) SearchOrderHistories(ctx context.Context, request *orderhistorypb.SearchOrderHistoriesRequest) (*orderhistorypb.SearchOrderHistoriesResponse, error) {
	var filters *queries.OrderHistoryFilters

	if request.Filter != nil {
		filters = &queries.OrderHistoryFilters{
			Since:    request.Filter.Since.AsTime(),
			Keywords: request.Filter.Keywords,
			Status:   orderapi.FromOrderStateProto(request.Filter.Status),
		}
	}

	results, err := h.app.Queries.SearchOrderHistories.Handle(ctx, queries.SearchOrderHistories{
		ConsumerID: request.ConsumerID,
		Filter:     filters,
		Next:       request.Next,
		Limit:      int(request.Limit),
	})
	if err != nil {
		return nil, err
	}

	orders := make([]*orderhistorypb.OrderHistory, len(results.Orders))
	for _, order := range results.Orders {
		orders = append(orders, h.toOrderHistoryProto(order))
	}

	return &orderhistorypb.SearchOrderHistoriesResponse{
		Orders: orders,
		Next:   results.Next,
	}, nil
}

func (h rpcHandlers) GetOrderHistory(ctx context.Context, request *orderhistorypb.GetOrderHistoryRequest) (*orderhistorypb.GetOrderHistoryResponse, error) {
	orderHistory, err := h.app.Queries.GetOrderHistory.Handle(ctx, queries.GetOrderHistory{OrderID: request.OrderID})
	if err != nil {
		return nil, err
	}

	return &orderhistorypb.GetOrderHistoryResponse{
		Order: h.toOrderHistoryProto(orderHistory),
	}, nil
}

func (h rpcHandlers) toOrderHistoryProto(orderHistory *domain.OrderHistory) *orderhistorypb.OrderHistory {
	return &orderhistorypb.OrderHistory{
		OrderID:        orderHistory.OrderID,
		Status:         orderapi.ToOrderStateProto(orderHistory.Status),
		RestaurantID:   orderHistory.RestaurantID,
		RestaurantName: orderHistory.RestaurantName,
		CreatedAt:      timestamppb.New(orderHistory.CreatedAt),
	}
}
