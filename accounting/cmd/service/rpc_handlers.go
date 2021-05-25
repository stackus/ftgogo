package main

import (
	"context"

	"github.com/stackus/ftgogo/account/internal/application/commands"
	"github.com/stackus/ftgogo/account/internal/application/queries"
	"github.com/stackus/ftgogo/serviceapis/accountingapi/pb"
)

type rpcHandlers struct {
	app Application
	accountingpb.UnimplementedAccountingServiceServer
}

var _ accountingpb.AccountingServiceServer = (*rpcHandlers)(nil)

func newRpcHandlers(app Application) rpcHandlers {
	return rpcHandlers{app: app}
}

func (h rpcHandlers) GetAccount(ctx context.Context, request *accountingpb.GetAccountRequest) (*accountingpb.GetAccountResponse, error) {
	_, err := h.app.Queries.GetAccount.Handle(ctx, queries.GetAccount{AccountID: request.AccountID})
	if err != nil {
		return nil, err
	}

	return &accountingpb.GetAccountResponse{AccountID: request.AccountID}, nil
}

func (h rpcHandlers) DisableAccount(ctx context.Context, request *accountingpb.DisableAccountRequest) (*accountingpb.DisableAccountResponse, error) {
	err := h.app.Commands.DisableAccount.Handle(ctx, commands.DisableAccount{AccountID: request.AccountID})
	if err != nil {
		return nil, err
	}

	return &accountingpb.DisableAccountResponse{AccountID: request.AccountID}, nil
}

func (h rpcHandlers) EnableAccount(ctx context.Context, request *accountingpb.EnableAccountRequest) (*accountingpb.EnableAccountResponse, error) {
	err := h.app.Commands.EnableAccount.Handle(ctx, commands.EnableAccount{AccountID: request.AccountID})
	if err != nil {
		return nil, err
	}

	return &accountingpb.EnableAccountResponse{AccountID: request.AccountID}, nil
}
