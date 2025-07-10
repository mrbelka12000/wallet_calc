//go:generate mockgen --source=interfaces.go --destination=mocks.go --package=usecase
package usecase

import (
	"context"

	"gorm.io/gorm"

	domainmodel "github.com/mrbelka12000/wallet_calc/internal/domain_model"
	"github.com/mrbelka12000/wallet_calc/internal/entity"
	"github.com/mrbelka12000/wallet_calc/pkg/gorm/postgres"
)

type (
	userRepository interface {
		Create(db *gorm.DB, req entity.User) error
		Get(db *gorm.DB, req entity.User) (entity.User, error)
	}

	accountRepository interface {
		List(db *gorm.DB, req entity.Account) ([]entity.Account, error)
		Save(db *gorm.DB, req entity.Account) error
	}

	categoryRepository interface {
		Create(db *gorm.DB, req entity.Category) error
		List(db *gorm.DB) ([]entity.Category, error)
	}

	db interface {
		TxBegin(ctx context.Context) *postgres.Gorm
		WithCtx(ctx context.Context) *postgres.Gorm
	}

	baseParser interface {
		ParseStatement(ctx context.Context, fileName string, categories []string) ([]domainmodel.BaseTransaction, error)
	}

	categoryUsecase interface {
		List(ctx context.Context) ([]domainmodel.Category, error)
	}
)
