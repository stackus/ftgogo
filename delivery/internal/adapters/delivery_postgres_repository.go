package adapters

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/stackus/edat-pgx"

	"github.com/stackus/ftgogo/delivery/internal/domain"
)

const (
	findDeliverySQL   = "SELECT restaurant_id, courier_id, pickup_address, delivery_address, pickup_time, ready_by, status FROM %s WHERE id = $1"
	saveDeliverySQL   = "INSERT INTO %s (id, restaurant_id, courier_id, pickup_address, delivery_address, pickup_time, ready_by, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
	updateDeliverySQL = "UPDATE %s SET restaurant_id = $1, courier_id = $2, pickup_address = $3, delivery_address = $4, pickup_time = $5, ready_by = $6, status = $7 WHERE id = $8"
)

type DeliveryPostgresRepository struct {
	client edatpgx.Client
}

var DeliveriesTableName = "deliveries"

var _ DeliveryRepository = (*DeliveryPostgresRepository)(nil)

func NewDeliveryPostgresRepository(client edatpgx.Client) *DeliveryPostgresRepository {
	return &DeliveryPostgresRepository{
		client: client,
	}
}

func (r *DeliveryPostgresRepository) Find(ctx context.Context, deliveryID string) (*domain.Delivery, error) {
	row := r.client.QueryRow(ctx, fmt.Sprintf(findDeliverySQL, DeliveriesTableName), deliveryID)

	var pickupData []byte
	var deliveryData []byte

	delivery := &domain.Delivery{
		DeliveryID: deliveryID,
	}

	err := row.Scan(
		&delivery.RestaurantID, &delivery.AssignedCourierID,
		pickupData, deliveryData,
		&delivery.PickUpTime, &delivery.ReadyBy, &delivery.Status,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrDeliveryNotFound
		}
		return nil, err
	}

	err = json.Unmarshal(pickupData, &delivery.PickUpAddress)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(deliveryData, &delivery.DeliveryAddress)
	if err != nil {
		return nil, err
	}

	return delivery, nil
}

func (r *DeliveryPostgresRepository) Save(ctx context.Context, delivery *domain.Delivery) error {
	var err error
	var pickupData []byte
	var deliveryData []byte

	pickupData, err = json.Marshal(delivery.PickUpAddress)
	if err != nil {
		return err
	}

	deliveryData, err = json.Marshal(delivery.DeliveryAddress)
	if err != nil {
		return err
	}

	_, err = r.client.Exec(ctx, fmt.Sprintf(saveDeliverySQL, DeliveriesTableName), delivery.DeliveryID, delivery.RestaurantID, delivery.AssignedCourierID,
		pickupData, deliveryData,
		delivery.PickUpTime, delivery.ReadyBy, delivery.Status.String(),
	)

	return err
}

func (r *DeliveryPostgresRepository) Update(ctx context.Context, deliveryID string, delivery *domain.Delivery) error {
	var err error
	var pickupData []byte
	var deliveryData []byte

	pickupData, err = json.Marshal(delivery.PickUpAddress)
	if err != nil {
		return err
	}

	deliveryData, err = json.Marshal(delivery.DeliveryAddress)
	if err != nil {
		return err
	}

	_, err = r.client.Exec(ctx, fmt.Sprintf(updateDeliverySQL, DeliveriesTableName), delivery.RestaurantID, delivery.AssignedCourierID,
		pickupData, deliveryData,
		delivery.PickUpTime, delivery.ReadyBy, delivery.Status.String(), deliveryID,
	)

	return err
}
