package adapters

import (
	"context"

	"github.com/stackus/ftgogo/serviceapis/accountingapi/pb"
	"github.com/stackus/ftgogo/store-web/internal/domain"
)

type AccountingGrpcRepository struct {
	client accountingpb.AccountingServiceClient
}

var _ AccountingRepository = (*AccountingGrpcRepository)(nil)

func NewAccountingGrpcRepository(client accountingpb.AccountingServiceClient) *AccountingGrpcRepository {
	return &AccountingGrpcRepository{client: client}
}

func (r AccountingGrpcRepository) Find(ctx context.Context, findAccount FindAccount) (*domain.Account, error) {
	resp, err := r.client.GetAccount(ctx, &accountingpb.GetAccountRequest{AccountID: findAccount.AccountID})
	if err != nil {
		return nil, err
	}

	return &domain.Account{
		AccountID: resp.AccountID,
		Enabled:   false,
	}, nil
}

func (r AccountingGrpcRepository) Enable(ctx context.Context, enableAccount EnableAccount) error {
	_, err := r.client.EnableAccount(ctx, &accountingpb.EnableAccountRequest{AccountID: enableAccount.AccountID})
	return err
}

func (r AccountingGrpcRepository) Disable(ctx context.Context, disableAccount DisableAccount) error {
	_, err := r.client.DisableAccount(ctx, &accountingpb.DisableAccountRequest{AccountID: disableAccount.AccountID})
	return err
}
