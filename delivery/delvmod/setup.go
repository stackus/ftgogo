package delvmod

import (
	"github.com/stackus/ftgogo/delivery/internal/adapters"
	"github.com/stackus/ftgogo/delivery/internal/application"
	"github.com/stackus/ftgogo/delivery/internal/handlers"
	"shared-go/applications"
)

func Setup(svc *applications.Monolith) error {
	// Driven
	adapters.CouriersTableName = "delivery.couriers"
	courierRepo := adapters.NewCourierPostgresRepository(svc.PgConn)
	adapters.DeliveriesTableName = "delivery.deliveries"
	deliveryRepo := adapters.NewDeliveryPostgresRepository(svc.PgConn)
	adapters.RestaurantsTableName = "delivery.restaurants"
	restaurantRepo := adapters.NewRestaurantPostgresRepository(svc.PgConn)

	app := application.NewServiceApplication(courierRepo, deliveryRepo, restaurantRepo)

	// Drivers
	handlers.NewRestaurantEventHandlers(app).Mount(svc.Subscriber)
	handlers.NewOrderEventHandlers(app).Mount(svc.Subscriber)
	handlers.NewTicketEventHandlers(app).Mount(svc.Subscriber)
	handlers.NewRpcHandlers(app).Mount(svc.RpcServer)

	return nil
}
