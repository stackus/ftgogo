package adapters

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/stackus/edat-pgx"

	"github.com/stackus/ftgogo/delivery/internal/domain"
)

const (
	findCourierSQL        = "SELECT plan, available FROM %s WHERE id = $1"
	findFirstAvailableSQL = "SELECT id, plan, available FROM %s WHERE available ORDER BY modified_at DESC LIMIT 1"
	saveCourierSQL        = "INSERT INTO %s (id, plan, available, modified_at) VALUES ($1, $2, $3, CURRENT_TIMESTAMP)"
	updateCourierSQL      = "UPDATE %s SET plan = $1, address = $2, modified_at = CURRENT_TIMESTAMP WHERE id = $3"
)

type CourierPostgresRepository struct {
	client edatpgx.Client
}

var CouriersTableName = "couriers"

var _ CourierRepository = (*CourierPostgresRepository)(nil)

func NewCourierPostgresRepository(client edatpgx.Client) *CourierPostgresRepository {
	return &CourierPostgresRepository{
		client: client,
	}
}

func (r *CourierPostgresRepository) Find(ctx context.Context, courierID string) (*domain.Courier, error) {
	row := r.client.QueryRow(ctx, fmt.Sprintf(findCourierSQL, CouriersTableName), courierID)

	var planData []byte

	c := &domain.Courier{
		CourierID: courierID,
	}

	err := row.Scan(&planData, &c.Available)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrCourierNotFound
		}
		return nil, err
	}

	err = json.Unmarshal(planData, &c.Plan)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (r *CourierPostgresRepository) FindOrCreate(ctx context.Context, courierID string) (*domain.Courier, error) {
	courier, err := r.Find(ctx, courierID)
	if err != nil {
		if errors.Is(err, domain.ErrCourierNotFound) {
			courier = &domain.Courier{
				CourierID: courierID,
				Plan:      domain.Plan{},
				Available: true,
			}
			err = r.Save(ctx, courier)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return courier, nil
}

func (r *CourierPostgresRepository) FindFirstAvailable(ctx context.Context) (*domain.Courier, error) {
	row := r.client.QueryRow(ctx, fmt.Sprintf(findFirstAvailableSQL, CouriersTableName))

	var planData []byte

	c := &domain.Courier{}

	err := row.Scan(&c.CourierID, &planData, &c.Available)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// for the demo only; in the real world you can't count on instant courier (just add water!)
			return r.FindOrCreate(ctx, uuid.New().String())
		}
		return nil, err
	}

	err = json.Unmarshal(planData, &c.Plan)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (r *CourierPostgresRepository) Save(ctx context.Context, courier *domain.Courier) error {
	planData, err := json.Marshal(courier.Plan)
	if err != nil {
		return err
	}

	_, err = r.client.Exec(ctx, fmt.Sprintf(saveCourierSQL, CouriersTableName), courier.CourierID, planData, courier.Available)

	return err
}

func (r *CourierPostgresRepository) Update(ctx context.Context, courierID string, courier *domain.Courier) error {
	planData, err := json.Marshal(courier.Plan)
	if err != nil {
		return err
	}

	_, err = r.client.Exec(ctx, fmt.Sprintf(updateCourierSQL, CouriersTableName), planData, courier.Available, courierID)

	return err
}
