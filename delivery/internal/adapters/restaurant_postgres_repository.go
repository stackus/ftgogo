package adapters

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/stackus/edat-pgx"

	"github.com/stackus/ftgogo/delivery/internal/application/ports"
	"github.com/stackus/ftgogo/delivery/internal/domain"
)

const (
	findRestaurantSQL   = "SELECT name, address FROM %s WHERE id = $1"
	saveRestaurantSQL   = "INSERT INTO %s (id, name, address) VALUES ($1, $2, $3)"
	updateRestaurantSQL = "UPDATE %s SET name = $1, address = $2 WHERE id = $3"
)

type RestaurantPostgresRepository struct {
	client edatpgx.Client
}

var RestaurantsTableName = "restaurants"

var _ ports.RestaurantRepository = (*RestaurantPostgresRepository)(nil)

func NewRestaurantPostgresRepository(client edatpgx.Client) *RestaurantPostgresRepository {
	return &RestaurantPostgresRepository{
		client: client,
	}
}

func (r *RestaurantPostgresRepository) Find(ctx context.Context, restaurantID string) (*domain.Restaurant, error) {
	row := r.client.QueryRow(ctx, fmt.Sprintf(findRestaurantSQL, RestaurantsTableName), restaurantID)

	var name string
	var addressData []byte

	err := row.Scan(&name, &addressData)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrRestaurantNotFound
		}
		return nil, err
	}

	restaurant := &domain.Restaurant{
		RestaurantID: restaurantID,
		Name:         name,
	}

	err = json.Unmarshal(addressData, &restaurant.Address)
	if err != nil {
		return nil, err
	}

	return restaurant, nil
}

func (r *RestaurantPostgresRepository) Save(ctx context.Context, restaurant *domain.Restaurant) error {
	addressData, err := json.Marshal(restaurant.Address)
	if err != nil {
		return err
	}

	_, err = r.client.Exec(ctx, fmt.Sprintf(saveRestaurantSQL, RestaurantsTableName), restaurant.RestaurantID, restaurant.Name, addressData)

	return err
}

func (r *RestaurantPostgresRepository) Update(ctx context.Context, restaurantID string, restaurant *domain.Restaurant) error {
	addressData, err := json.Marshal(restaurant.Address)
	if err != nil {
		return err
	}

	_, err = r.client.Exec(ctx, fmt.Sprintf(updateRestaurantSQL, RestaurantsTableName), restaurant.Name, addressData, restaurantID)

	return err
}
