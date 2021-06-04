package main

import (
	"github.com/stackus/ftgogo/accounting/internal/adapters"
	"github.com/stackus/ftgogo/accounting/internal/application"
	"github.com/stackus/ftgogo/accounting/internal/application/commands"
	"github.com/stackus/ftgogo/accounting/internal/application/queries"
	"github.com/stackus/ftgogo/accounting/internal/domain"
	"github.com/stackus/ftgogo/accounting/internal/handlers"
	"github.com/stackus/ftgogo/serviceapis"
	"shared-go/applications"
)

func main() {
	svc := applications.NewService(initService)
	if err := svc.Execute(); err != nil {
		panic(err)
	}
}

func initService(svc *applications.Service) error {
	serviceapis.RegisterTypes()
	domain.RegisterTypes()

	// Driven
	accountRepo := adapters.NewAccountAggregateRootRepository(svc.AggregateStore)

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

	return nil
}
