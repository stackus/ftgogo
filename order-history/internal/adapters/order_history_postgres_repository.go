package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/stackus/edat-pgx"

	"github.com/stackus/ftgogo/order-history/internal/domain"
	"serviceapis/orderapi"
)

const (
	findOrderHistorySQL   = "SELECT consumer_id, restaurant_id, restaurant_name, line_items, order_total, status, keywords, created_at FROM orders WHERE id = $1"
	saveOrderHistorySQL   = "INSERT INTO orders (id, consumer_id, restaurant_id, restaurant_name, line_items, order_total, status, keywords, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"
	updateOrderStatusSQL  = "UPDATE orders SET status = $1 WHERE id = $2"
	updateOrderHistorySQL = "UPDATE orders SET consumer_id = $1, restaurant_id = $2, restaurant_name = $3, line_items = $4, order_total = $5, status = $6, keywords = $7 WHERE id = $8"
)

type OrderHistoryPostgresRepository struct {
	client edatpgx.Client
}

var _ domain.OrderHistoryRepository = (*OrderHistoryPostgresRepository)(nil)

func NewOrderHistoryPostgresRepository(client edatpgx.Client) *OrderHistoryPostgresRepository {
	return &OrderHistoryPostgresRepository{
		client: client,
	}
}

func (r *OrderHistoryPostgresRepository) FindConsumerOrders(ctx context.Context, consumerID string, filters domain.OrderHistoryFilters) ([]*domain.OrderHistory, string, error) {
	type cursor struct {
		CreatedAt time.Time `json:"created_at"`
		ID        string    `json:"id"`
	}

	query := "SELECT id, restaurant_id, restaurant_name, line_items, order_total, status, keywords, created_at FROM orders WHERE consumer_id = $1"

	paramCount := 1
	params := []interface{}{consumerID}

	if !filters.Since.IsZero() {
		query += fmt.Sprintf(" AND created_at > $%d", paramCount+1)
		params = append(params, filters.Since)
		paramCount += 1
	}

	if filters.Status != orderapi.UnknownOrderState {
		query += fmt.Sprintf(" AND status = $%d", paramCount+1)
		params = append(params, int(filters.Status))
		paramCount += 1
	}

	if len(filters.Keywords) > 0 {
		query += fmt.Sprintf(" AND keywords @> $%d", paramCount+1)
		params = append(params, filters.Keywords)
		paramCount += 1
	}

	if filters.Next != "" {
		seek := cursor{}
		err := json.Unmarshal([]byte(filters.Next), &seek)
		if err != nil {
			return nil, "", err
		}

		query += fmt.Sprintf(" AND (created_at, id) < ($%d, $%d)", paramCount+1, paramCount+2)
		paramCount += 2
		params = append(params, seek.CreatedAt, seek.ID)
	}

	query += " ORDER BY created_at DESC, id DESC"
	query += fmt.Sprintf(" LIMIT %d", filters.Limit)

	rows, err := r.client.Query(ctx, query, params...)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	orders := []*domain.OrderHistory{}

	for rows.Next() {
		lineItemData := []byte{}
		order := &domain.OrderHistory{
			ConsumerID: consumerID,
		}

		err := rows.Scan(
			&order.OrderID, &order.RestaurantID, &order.RestaurantName,
			&lineItemData, &order.OrderTotal,
			&order.Status, &order.Keywords, &order.CreatedAt,
		)
		if err != nil {
			return nil, "", err
		}
		err = json.Unmarshal(lineItemData, &order.LineItems)
		if err != nil {
			return nil, "", err
		}

		orders = append(orders, order)
	}

	if len(orders) == 0 {
		return orders, "", nil
	}

	lastOrder := orders[len(orders)-1]

	seek := cursor{
		ID:        lastOrder.OrderID,
		CreatedAt: lastOrder.CreatedAt,
	}

	var next []byte
	next, err = json.Marshal(seek)
	if err != nil {
		return nil, "", err
	}

	return orders, string(next), nil
}

func (r *OrderHistoryPostgresRepository) Find(ctx context.Context, orderHistoryID string) (*domain.OrderHistory, error) {
	row := r.client.QueryRow(ctx, findOrderHistorySQL, orderHistoryID)

	lineItemData := []byte{}
	order := &domain.OrderHistory{
		OrderID: orderHistoryID,
	}

	err := row.Scan(
		&order.ConsumerID, &order.RestaurantID, &order.RestaurantName,
		&lineItemData, &order.OrderTotal,
		&order.Status, &order.Keywords, &order.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(lineItemData, &order.LineItems)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (r *OrderHistoryPostgresRepository) Save(ctx context.Context, orderHistory *domain.OrderHistory) error {
	lineItemData, err := json.Marshal(orderHistory.LineItems)
	if err != nil {
		return err
	}

	_, err = r.client.Exec(ctx, saveOrderHistorySQL,
		orderHistory.OrderID, orderHistory.ConsumerID, orderHistory.RestaurantID,
		orderHistory.RestaurantName, lineItemData, orderHistory.OrderTotal,
		int(orderHistory.Status), orderHistory.Keywords, orderHistory.CreatedAt,
	)
	return err
}

func (r *OrderHistoryPostgresRepository) UpdateStatus(ctx context.Context, orderHistoryID string, status orderapi.OrderState) error {
	_, err := r.client.Exec(ctx, updateOrderStatusSQL, int(status), orderHistoryID)

	return err
}

func (r *OrderHistoryPostgresRepository) Update(ctx context.Context, orderHistoryID string, orderHistory *domain.OrderHistory) error {
	lineItemData, err := json.Marshal(orderHistory.LineItems)
	if err != nil {
		return err
	}

	_, err = r.client.Exec(ctx, updateOrderHistorySQL,
		orderHistory.ConsumerID, orderHistory.RestaurantID, orderHistory.RestaurantName,
		lineItemData, orderHistory.OrderTotal,
		int(orderHistory.Status), orderHistory.Keywords,
		orderHistoryID,
	)
	return err
}
