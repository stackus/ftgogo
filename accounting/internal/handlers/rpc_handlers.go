package handlers

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/stackus/ftgogo/accounting/internal/application"
	"github.com/stackus/ftgogo/accounting/internal/application/commands"
	"github.com/stackus/ftgogo/accounting/internal/application/queries"
	"github.com/stackus/ftgogo/serviceapis/accountingapi/pb"
)

type RpcHandlers struct {
	app application.Service
	accountingpb.UnimplementedAccountingServiceServer
}

var _ accountingpb.AccountingServiceServer = (*RpcHandlers)(nil)

func NewRpcHandlers(app application.Service) RpcHandlers {
	return RpcHandlers{app: app}
}

func (h RpcHandlers) Mount(registrar grpc.ServiceRegistrar) {
	accountingpb.RegisterAccountingServiceServer(registrar, h)
}

func (h RpcHandlers) GetAccount(ctx context.Context, request *accountingpb.GetAccountRequest) (*accountingpb.GetAccountResponse, error) {
	account, err := h.app.Queries.GetAccount.Handle(ctx, queries.GetAccount{AccountID: request.AccountID})
	if err != nil {
		return nil, err
	}

	return &accountingpb.GetAccountResponse{
		AccountID: request.AccountID,
		Enabled:   account.Enabled,
	}, nil
}

func (h RpcHandlers) DisableAccount(ctx context.Context, request *accountingpb.DisableAccountRequest) (*emptypb.Empty, error) {
	err := h.app.Commands.DisableAccount.Handle(ctx, commands.DisableAccount{AccountID: request.AccountID})
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (h RpcHandlers) EnableAccount(ctx context.Context, request *accountingpb.EnableAccountRequest) (*emptypb.Empty, error) {
	err := h.app.Commands.EnableAccount.Handle(ctx, commands.EnableAccount{AccountID: request.AccountID})
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
