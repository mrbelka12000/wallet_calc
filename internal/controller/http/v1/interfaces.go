package v1

import (
	"context"

	domainmodel "github.com/mrbelka12000/wallet_calc/internal/domain_model"
	"github.com/mrbelka12000/wallet_calc/internal/dto"
)

type (
	userUsecase interface {
		Register(ctx context.Context, user domainmodel.User) error
		Login(ctx context.Context, user domainmodel.User) (dto.User, error)
		Get(ctx context.Context, req domainmodel.User) (dto.User, error)
	}

	transactionUsecase interface {
	}
)
