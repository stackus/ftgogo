package consmod

import (
	edatpgx "github.com/stackus/edat-pgx"
	"github.com/stackus/edat/outbox"

	"github.com/stackus/ftgogo/consumer/internal/adapters"
	"github.com/stackus/ftgogo/consumer/internal/application"
	"github.com/stackus/ftgogo/consumer/internal/application/commands"
	"github.com/stackus/ftgogo/consumer/internal/application/queries"
	"github.com/stackus/ftgogo/consumer/internal/domain"
	"github.com/stackus/ftgogo/consumer/internal/handlers"
	"shared-go/applications"
)

func Setup(svc *applications.Monolith) error {
	domain.RegisterTypes()

	// Infrastructure
	aggregateStore := edatpgx.NewSnapshotStore(
		svc.PgConn,
		edatpgx.WithSnapshotStoreTableName("consumer.snapshots"),
	)(edatpgx.NewEventStore(
		svc.PgConn,
		edatpgx.WithEventStoreTableName("consumer.events"),
	))
	messageStore := edatpgx.NewMessageStore(svc.PgConn, edatpgx.WithMessageStoreTableName("consumer.messages"))

	// Driven
	consumerRepo := adapters.NewConsumerRepositoryPublisherMiddleware(
		adapters.NewConsumerAggregateRootRepository(aggregateStore),
		adapters.NewConsumerEntityEventPublisher(svc.Publisher),
	)

	app := application.Service{
		Commands: application.Commands{
			RegisterConsumer:        commands.NewRegisterConsumerHandler(consumerRepo),
			UpdateConsumer:          commands.NewUpdateConsumerHandler(consumerRepo),
			ValidateOrderByConsumer: commands.NewValidateOrderByConsumerHandler(consumerRepo),
			AddAddress:              commands.NewAddAddressHandler(consumerRepo),
			UpdateAddress:           commands.NewUpdateAddressHandler(consumerRepo),
			RemoveAddress:           commands.NewRemoveAddressHandler(consumerRepo),
		},
		Queries: application.Queries{
			GetConsumer: queries.NewGetConsumerHandler(consumerRepo),
			GetAddress:  queries.NewGetAddressHandler(consumerRepo),
		},
	}

	// Drivers
	handlers.NewCommandHandlers(app).Mount(svc.Subscriber, svc.Publisher)
	handlers.NewRpcHandlers(app).Mount(svc.RpcServer)
	svc.Processors = append(svc.Processors, outbox.NewPollingProcessor(messageStore, svc.CDCPublisher))

	return nil
}
