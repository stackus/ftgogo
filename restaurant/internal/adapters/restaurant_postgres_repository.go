package adapters

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/stackus/edat-pgx"

	"github.com/stackus/ftgogo/restaurant/internal/domain"
)

const (
	findRestaurantSQL   = "SELECT name, address, menu FROM %s WHERE id = $1"
	saveRestaurantSQL   = "INSERT INTO %s (id, name, address, menu) VALUES ($1, $2, $3, $4)"
	updateRestaurantSQL = "UPDATE %s SET name = $1, address = $2, menu = $3 WHERE id = $4"
)

type RestaurantPostgresRepository struct {
	client edatpgx.Client
}

var RestaurantsTableName = "restaurants"

var _ RestaurantRepository = (*RestaurantPostgresRepository)(nil)

func NewRestaurantPostgresRepository(client edatpgx.Client) *RestaurantPostgresRepository {
	return &RestaurantPostgresRepository{
		client: client,
	}
}

func (s *RestaurantPostgresRepository) Find(ctx context.Context, restaurantID string) (*domain.Restaurant, error) {
	row := s.client.QueryRow(ctx, fmt.Sprintf(findRestaurantSQL, RestaurantsTableName), restaurantID)

	var name string
	var addressData []byte
	var menuData []byte

	err := row.Scan(&name, &addressData, &menuData)
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

	err = json.Unmarshal(menuData, &r.MenuItems)
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

	menuData, err := json.Marshal(r.MenuItems)
	if err != nil {
		return err
	}

	_, err = s.client.Exec(ctx, fmt.Sprintf(saveRestaurantSQL, RestaurantsTableName), r.RestaurantID, r.Name, addressData, menuData)

	return err
}

func (s *RestaurantPostgresRepository) Update(ctx context.Context, restaurantID string, r *domain.Restaurant) error {
	addressData, err := json.Marshal(r.Address)
	if err != nil {
		return err
	}

	menuData, err := json.Marshal(r.MenuItems)
	if err != nil {
		return err
	}

	_, err = s.client.Exec(ctx, fmt.Sprintf(updateRestaurantSQL, RestaurantsTableName), r.Name, addressData, menuData, restaurantID)

	return err
}
