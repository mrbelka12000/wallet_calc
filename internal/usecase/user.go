package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/mrbelka12000/wallet_calc/internal/converter"
	domainmodel "github.com/mrbelka12000/wallet_calc/internal/domain_model"
	"github.com/mrbelka12000/wallet_calc/internal/dto"
	"github.com/mrbelka12000/wallet_calc/internal/entity"
	"github.com/mrbelka12000/wallet_calc/pkg/validator"
)

type (
	UserUsecase struct {
		repo userRepository
		db   db
	}
)

func NewUserUsecase(repo userRepository, db db) *UserUsecase {
	return &UserUsecase{
		repo: repo,
		db:   db,
	}
}

func (uc *UserUsecase) Register(ctx context.Context, req domainmodel.User) error {
	if err := validator.ValidateStruct(req); err != nil {
		return newClientError(err.Error())
	}

	passwordHash, err := generateHashFromPassword(req.Password)
	if err != nil {
		return err
	}

	userReq := entity.User{
		ID:        uuid.New(),
		Name:      req.Name,
		Email:     req.Email,
		Password:  string(passwordHash),
		CreatedAt: time.Now().UTC(),
	}

	tx := uc.db.TxBegin(ctx)
	if tx.Error != nil {
		return tx.Error
	}
	defer tx.Rollback()

	if _, err := uc.repo.Get(tx.DB, entity.User{Email: req.Email}); err == nil {
		return newClientError("Email already used")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if err := uc.repo.Create(tx.DB, userReq); err != nil {
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (uc *UserUsecase) Login(ctx context.Context, req domainmodel.User) (dto.User, error) {
	userReq := entity.User{
		Email: req.Email,
	}

	entityUser, err := uc.repo.Get(uc.db.WithCtx(ctx).DB, userReq)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.User{}, newClientError("User not found")
		}

		return dto.User{}, err
	}

	if !isValidPassword(entityUser.Password, req.Password) {
		return dto.User{}, newClientError("Email or Password is incorrect")
	}

	return converter.ConvertFromEntityUserToDTO(entityUser, false), nil
}

func (uc *UserUsecase) Get(ctx context.Context, req domainmodel.User) (dto.User, error) {
	entityUser, err := uc.repo.Get(uc.db.WithCtx(ctx).DB, entity.ConvertFromDMUser(req))
	if err != nil {
		return dto.User{}, err
	}

	return converter.ConvertFromEntityUserToDTO(entityUser, false), nil
}

func generateHashFromPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func isValidPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
