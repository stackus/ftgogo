package restmod

import (
	edatpgx "github.com/stackus/edat-pgx"
	"github.com/stackus/edat/msg"
	"github.com/stackus/edat/outbox"

	"github.com/stackus/ftgogo/restaurant/internal/adapters"
	"github.com/stackus/ftgogo/restaurant/internal/application"
	"github.com/stackus/ftgogo/restaurant/internal/handlers"
	"github.com/stackus/ftgogo/serviceapis"
	"shared-go/applications"
)

func Setup(svc *applications.Monolith) error {
	serviceapis.RegisterTypes()

	// Infrastructure
	messageStore := edatpgx.NewMessageStore(svc.CDCPgConn, edatpgx.WithMessageStoreTableName("restaurant.messages"))
	publisher := msg.NewPublisher(messageStore)
	svc.Publishers = append(svc.Publishers, publisher)
	svc.Processors = append(svc.Processors, outbox.NewPollingProcessor(messageStore, svc.CDCPublisher))

	// Driven
	adapters.RestaurantsTableName = "restaurant.restaurants"
	restaurantRepo := adapters.NewRestaurantPostgresPublisherMiddleware(
		adapters.NewRestaurantPostgresRepository(svc.PgConn),
		adapters.NewRestaurantEntityEventPublisher(publisher),
	)

	app := application.NewServiceApplication(restaurantRepo)

	// Drivers
	handlers.NewRpcHandlers(app).Mount(svc.RpcServer)

	return nil
}
