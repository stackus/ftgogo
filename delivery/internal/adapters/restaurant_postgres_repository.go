package adapters

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/stackus/edat-pgx"

	"github.com/stackus/ftgogo/delivery/internal/domain"
)

const (
	findRestaurantSQL   = "SELECT name, address FROM restaurants WHERE id = $1"
	saveRestaurantSQL   = "INSERT INTO restaurants (id, name, address) VALUES ($1, $2, $3)"
	updateRestaurantSQL = "UPDATE restaurants SET name = $1, address = $2 WHERE id = $3"
)

type RestaurantPostgresRepository struct {
	client edatpgx.Client
}

var _ domain.RestaurantRepository = (*RestaurantPostgresRepository)(nil)

func NewRestaurantPostgresRepository(client edatpgx.Client) *RestaurantPostgresRepository {
	return &RestaurantPostgresRepository{
		client: client,
	}
}

func (s *RestaurantPostgresRepository) Find(ctx context.Context, restaurantID string) (*domain.Restaurant, error) {
	row := s.client.QueryRow(ctx, findRestaurantSQL, restaurantID)

	var name string
	var addressData []byte

	err := row.Scan(&name, &addressData)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrRestaurantNotFound
		}
		return nil, err
	}

	r := &domain.Restaurant{
		RestaurantID: restaurantID,
		Name:         name,
	}

	err = json.Unmarshal(addressData, &r.Address)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (s *RestaurantPostgresRepository) Save(ctx context.Context, r *domain.Restaurant) error {
	addressData, err := json.Marshal(r.Address)
	if err != nil {
		return err
	}

	_, err = s.client.Exec(ctx, saveRestaurantSQL, r.RestaurantID, r.Name, addressData)

	return err
}

func (s *RestaurantPostgresRepository) Update(ctx context.Context, restaurantID string, r *domain.Restaurant) error {
	addressData, err := json.Marshal(r.Address)
	if err != nil {
		return err
	}

	_, err = s.client.Exec(ctx, updateRestaurantSQL, r.Name, addressData, restaurantID)

	return err
}
