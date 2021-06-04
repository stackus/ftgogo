package acctmod

import (
	edatpgx "github.com/stackus/edat-pgx"
	"github.com/stackus/edat/outbox"

	"github.com/stackus/ftgogo/accounting/internal/adapters"
	"github.com/stackus/ftgogo/accounting/internal/application"
	"github.com/stackus/ftgogo/accounting/internal/application/commands"
	"github.com/stackus/ftgogo/accounting/internal/application/queries"
	"github.com/stackus/ftgogo/accounting/internal/domain"
	"github.com/stackus/ftgogo/accounting/internal/handlers"
	"shared-go/applications"
)

func Setup(svc *applications.Monolith) error {
	domain.RegisterTypes()

	// Infrastructure
	aggregateStore := edatpgx.NewSnapshotStore(
		svc.PgConn,
		edatpgx.WithSnapshotStoreTableName("accounting.snapshots"),
	)(edatpgx.NewEventStore(
		svc.PgConn,
		edatpgx.WithEventStoreTableName("accounting.events"),
	))
	messageStore := edatpgx.NewMessageStore(svc.PgConn, edatpgx.WithMessageStoreTableName("accounting.messages"))

	// Driven
	accountRepo := adapters.NewAccountAggregateRootRepository(aggregateStore)

	app := application.Service{
		Commands: application.Commands{
			AuthorizeOrder:        commands.NewAuthorizeOrderHandler(accountRepo),
			ReverseAuthorizeOrder: commands.NewReverseAuthorizeOrderHandler(accountRepo),
			ReviseAuthorizeOrder:  commands.NewReviseAuthorizeOrderHandler(accountRepo),
			CreateAccount:         commands.NewCreateAccountHandler(accountRepo),
			DisableAccount:        commands.NewDisableAccountHandler(accountRepo),
			EnableAccount:         commands.NewEnableAccountHandler(accountRepo),
		},
		Queries: application.Queries{
			GetAccount: queries.NewGetAccountHandler(accountRepo),
		},
	}

	// Drivers
	handlers.NewCommandHandlers(app).Mount(svc.Subscriber, svc.Publisher)
	handlers.NewConsumerEventHandlers(app).Mount(svc.Subscriber)
	handlers.NewRpcHandlers(app).Mount(svc.RpcServer)
	svc.Processors = append(svc.Processors, outbox.NewPollingProcessor(messageStore, svc.CDCPublisher))

	return nil
}
