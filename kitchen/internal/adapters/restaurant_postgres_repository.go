package adapters

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/stackus/edat-pgx"

	"github.com/stackus/ftgogo/kitchen/internal/domain"
)

const (
	findRestaurantSQL   = "SELECT name, menu FROM %s WHERE id = $1"
	saveRestaurantSQL   = "INSERT INTO %s (id, name, menu) VALUES ($1, $2, $3)"
	updateRestaurantSQL = "UPDATE %s SET name = $1, menu = $2 WHERE id = $3"
)

type RestaurantPostgresRepository struct {
	client edatpgx.Client
}

var RestaurantsTableName = "restaurants"

func NewRestaurantPostgresRepository(client edatpgx.Client) *RestaurantPostgresRepository {
	return &RestaurantPostgresRepository{
		client: client,
	}
}

func (s *RestaurantPostgresRepository) Find(ctx context.Context, restaurantID string) (*domain.Restaurant, error) {
	row := s.client.QueryRow(ctx, fmt.Sprintf(findRestaurantSQL, RestaurantsTableName), restaurantID)

	var name string
	var data []byte

	err := row.Scan(&name, &data)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrRestaurantNotFound
		}
		return nil, err
	}

	restaurant := &domain.Restaurant{
		RestaurantID: restaurantID,
		Name:         name,
	}

	err = json.Unmarshal(data, &restaurant.MenuItems)
	if err != nil {
		return nil, err
	}

	return restaurant, nil
}

func (s *RestaurantPostgresRepository) Save(ctx context.Context, restaurant *domain.Restaurant) error {
	menuItemData, err := json.Marshal(restaurant.MenuItems)
	if err != nil {
		return err
	}

	_, err = s.client.Exec(ctx, fmt.Sprintf(saveRestaurantSQL, RestaurantsTableName), restaurant.RestaurantID, restaurant.Name, menuItemData)
	if err != nil {
		return err
	}

	return nil
}

func (s *RestaurantPostgresRepository) Update(ctx context.Context, restaurantID string, restaurant *domain.Restaurant) error {
	menuItemData, err := json.Marshal(restaurant.MenuItems)
	if err != nil {
		return err
	}

	_, err = s.client.Exec(ctx, fmt.Sprintf(updateRestaurantSQL, RestaurantsTableName), restaurant.Name, menuItemData, restaurantID)
	if err != nil {
		return err
	}

	return nil
}
