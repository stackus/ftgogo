package restmod

import (
	edatpgx "github.com/stackus/edat-pgx"
	"github.com/stackus/edat/outbox"

	"github.com/stackus/ftgogo/restaurant/internal/adapters"
	"github.com/stackus/ftgogo/restaurant/internal/application"
	"github.com/stackus/ftgogo/restaurant/internal/application/commands"
	"github.com/stackus/ftgogo/restaurant/internal/application/queries"
	"github.com/stackus/ftgogo/restaurant/internal/handlers"
	"github.com/stackus/ftgogo/serviceapis"
	"shared-go/applications"
)

func Setup(svc *applications.Monolith) error {
	serviceapis.RegisterTypes()

	// Driven
	adapters.RestaurantsTableName = "restaurant.restaurants"
	restaurantRepo := adapters.NewRestaurantPostgresPublisherMiddleware(
		adapters.NewRestaurantPostgresRepository(svc.PgConn),
		adapters.NewRestaurantEntityEventPublisher(svc.Publisher),
	)
	messageStore := edatpgx.NewMessageStore(svc.PgConn, edatpgx.WithMessageStoreTableName("restaurant.messages"))

	app := application.Application{
		Commands: application.Commands{
			CreateRestaurant: commands.NewCreateRestaurantHandler(restaurantRepo),
		},
		Queries: application.Queries{
			GetRestaurant: queries.NewGetRestaurantHandler(restaurantRepo),
		},
	}

	// Drivers
	handlers.NewRpcHandlers(app).Mount(svc.RpcServer)
	svc.Processors = append(svc.Processors, outbox.NewPollingProcessor(messageStore, svc.CDCPublisher))

	return nil
}
