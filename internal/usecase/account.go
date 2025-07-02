package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/mrbelka12000/wallet_calc/internal/converter"
	domainmodel "github.com/mrbelka12000/wallet_calc/internal/domain_model"
	"github.com/mrbelka12000/wallet_calc/internal/dto"
	"github.com/mrbelka12000/wallet_calc/internal/entity"
	"github.com/mrbelka12000/wallet_calc/pkg/validator"
)

type (
	AccountUsecase struct {
		repo accountRepository
		db   db
	}
)

func NewAccountUsecase(repo accountRepository, db db) *AccountUsecase {
	return &AccountUsecase{
		repo: repo,
		db:   db,
	}
}

func (u *AccountUsecase) Create(ctx context.Context, req domainmodel.Account) error {
	if err := validator.ValidateStruct(req); err != nil {
		return newClientError(err.Error())
	}

	req.ID = uuid.New()

	accountReq := entity.Account{
		ID:        req.ID,
		UserID:    req.UserID,
		Name:      req.Name,
		Balance:   req.Balance,
		CreatedAt: time.Now().UTC(),
	}

	return u.repo.Save(u.db.WithCtx(ctx).DB, accountReq)
}

func (u *AccountUsecase) List(ctx context.Context, req domainmodel.Account) ([]dto.AccountResp, error) {
	accountReq := entity.Account{
		UserID: req.UserID,
		Name:   req.Name,
	}

	accounts, err := u.repo.List(u.db.WithCtx(ctx).DB, accountReq)
	if err != nil {
		return nil, err
	}

	return converter.ConvertFromEntityAccountsToDTO(accounts), nil
}

func (u *AccountUsecase) Update(ctx context.Context, req domainmodel.Account) error {
	if err := validator.ValidateStruct(req); err != nil {
		return newClientError(err.Error())
	}

	accountReq := entity.Account{
		ID:      req.ID,
		Name:    req.Name,
		Balance: req.Balance,
	}

	return u.repo.Save(u.db.WithCtx(ctx).DB, accountReq)
}
