package main

import (
	"github.com/stackus/edat/msg"
	"github.com/stackus/edat/saga"

	"github.com/stackus/ftgogo/account/internal/adapters"
	"github.com/stackus/ftgogo/account/internal/application/commands"
	"github.com/stackus/ftgogo/account/internal/application/queries"
	"github.com/stackus/ftgogo/account/internal/domain"
	"github.com/stackus/ftgogo/serviceapis"
	"github.com/stackus/ftgogo/serviceapis/accountingapi"
	"github.com/stackus/ftgogo/serviceapis/accountingapi/pb"
	"github.com/stackus/ftgogo/serviceapis/consumerapi"
	"shared-go/applications"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	AuthorizeOrder        commands.AuthorizeOrderHandler
	ReverseAuthorizeOrder commands.ReverseAuthorizeOrderHandler
	ReviseAuthorizeOrder  commands.ReviseAuthorizeOrderHandler
	CreateAccount         commands.CreateAccountHandler
	DisableAccount        commands.DisableAccountHandler
	EnableAccount         commands.EnableAccountHandler
}

type Queries struct {
	GetAccount queries.GetAccountHandler
}

func main() {
	svc := applications.NewService(initService)
	if err := svc.Execute(); err != nil {
		panic(err)
	}
}

func initService(svc *applications.Service) error {
	serviceapis.RegisterTypes()
	domain.RegisterTypes()

	accountRepo := adapters.NewAccountRepository(svc.AggregateStore)

	application := Application{
		Commands: Commands{
			AuthorizeOrder:        commands.NewAuthorizeOrderHandler(accountRepo),
			ReverseAuthorizeOrder: commands.NewReverseAuthorizeOrderHandler(accountRepo),
			ReviseAuthorizeOrder:  commands.NewReviseAuthorizeOrderHandler(accountRepo),
			CreateAccount:         commands.NewCreateAccountHandler(accountRepo),
			DisableAccount:        commands.NewDisableAccountHandler(accountRepo),
			EnableAccount:         commands.NewEnableAccountHandler(accountRepo),
		},
		Queries: Queries{
			GetAccount: queries.NewGetAccountHandler(accountRepo),
		},
	}

	cmdHandlers := NewCommandHandlers(application)
	svc.Subscriber.Subscribe(accountingapi.AccountingServiceCommandChannel, saga.NewCommandDispatcher(svc.Publisher).
		Handle(accountingapi.AuthorizeOrder{}, cmdHandlers.AuthorizeOrder).
		Handle(accountingapi.ReverseAuthorizeOrder{}, cmdHandlers.ReverseAuthorizeOrder).
		Handle(accountingapi.ReviseAuthorizeOrder{}, cmdHandlers.ReviseAuthorizeOrder))

	consumerEventHandlers := NewConsumerEventHandlers(application)
	svc.Subscriber.Subscribe(consumerapi.ConsumerAggregateChannel, msg.NewEntityEventDispatcher().
		Handle(consumerapi.ConsumerRegistered{}, consumerEventHandlers.ConsumerRegistered))

	accountingpb.RegisterAccountingServiceServer(svc.RpcServer, newRpcHandlers(application))

	return nil
}
