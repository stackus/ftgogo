package adapters

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/stackus/ftgogo/customer-web/internal/domain"
	"github.com/stackus/ftgogo/serviceapis/orderapi"
	"github.com/stackus/ftgogo/serviceapis/orderhistoryapi/pb"
)

type OrderHistoryGRPCRepository struct {
	client orderhistorypb.OrderHistoryServiceClient
}

var _ OrderHistoryRepository = (*OrderHistoryGRPCRepository)(nil)

func NewOrderHistoryGrpcClient(client orderhistorypb.OrderHistoryServiceClient) *OrderHistoryGRPCRepository {
	return &OrderHistoryGRPCRepository{client: client}
}

func (r OrderHistoryGRPCRepository) SearchOrders(ctx context.Context, searchOrders SearchOrders) (*SearchOrdersResult, error) {
	results, err := r.client.SearchOrderHistories(ctx, &orderhistorypb.SearchOrderHistoriesRequest{
		ConsumerID: searchOrders.ConsumerID,
		Filter:     r.toOrderHistoriesFilters(searchOrders.Filters),
		Next:       searchOrders.Next,
		Limit:      int64(searchOrders.Limit),
	})
	if err != nil {
		return nil, err
	}

	orders := make([]*domain.OrderHistory, 0, len(results.Orders))
	for _, order := range results.Orders {
		orders = append(orders, r.fromOrderHistoryProto(order))
	}

	return &SearchOrdersResult{
		Orders: orders,
		Next:   results.Next,
	}, nil
}

func (r OrderHistoryGRPCRepository) toOrderHistoriesFilters(filters *domain.SearchOrdersFilters) *orderhistorypb.SearchOrderHistoriesRequestFilters {
	return &orderhistorypb.SearchOrderHistoriesRequestFilters{
		Since:    timestamppb.New(filters.Since),
		Keywords: filters.Keywords,
		Status:   orderapi.ToOrderStateProto(filters.Status),
	}
}

func (r OrderHistoryGRPCRepository) fromOrderHistoryProto(order *orderhistorypb.OrderHistory) *domain.OrderHistory {
	return &domain.OrderHistory{
		OrderID:        order.OrderID,
		ConsumerID:     order.ConsumerID,
		RestaurantID:   order.RestaurantID,
		RestaurantName: order.RestaurantName,
		Status:         orderapi.FromOrderStateProto(order.Status),
		CreatedAt:      order.CreatedAt.AsTime(),
	}
}
