//go:generate mockgen --source=interfaces.go --destination=mocks.go --package=usecase
package usecase

import (
	"context"

	"gorm.io/gorm"

	"github.com/mrbelka12000/wallet_calc/internal/entity"
	"github.com/mrbelka12000/wallet_calc/pkg/gorm/postgres"
)

type (
	userRepository interface {
		Create(db *gorm.DB, req entity.User) error
		Get(db *gorm.DB, req entity.User) (entity.User, error)
	}

	db interface {
		TxBegin(ctx context.Context) *postgres.Gorm
		WithCtx(ctx context.Context) *postgres.Gorm
	}
)
