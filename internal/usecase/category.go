package usecase

import (
	"context"

	"github.com/google/uuid"

	"github.com/mrbelka12000/wallet_calc/internal/converter"
	domainmodel "github.com/mrbelka12000/wallet_calc/internal/domain_model"
	"github.com/mrbelka12000/wallet_calc/internal/entity"
	"github.com/mrbelka12000/wallet_calc/pkg/validator"
)

type (
	CategoryUsecase struct {
		repo categoryRepository
		db   db
	}
)

func NewCategoryUsecase(repo categoryRepository, db db) *CategoryUsecase {
	return &CategoryUsecase{
		repo: repo,
		db:   db,
	}
}

func (c *CategoryUsecase) Create(ctx context.Context, req domainmodel.Category) error {
	if err := validator.ValidateStruct(req); err != nil {
		return newClientError(err.Error())
	}
	req.ID = uuid.New()

	categoryReq := entity.Category{
		ID:   req.ID,
		Name: req.Name,
		Type: req.Type,
	}

	return c.repo.Create(c.db.WithCtx(ctx).DB, categoryReq)
}

func (c *CategoryUsecase) List(ctx context.Context) ([]domainmodel.Category, error) {
	categories, err := c.repo.List(c.db.WithCtx(ctx).DB)
	if err != nil {
		return nil, err
	}

	return converter.ConvertFromEntityCategoriesToDM(categories), nil
}
