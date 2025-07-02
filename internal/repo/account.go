package repo

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/mrbelka12000/wallet_calc/internal/entity"
)

type Account struct {
}

func NewAccount() *Account {
	return &Account{}
}

func (a *Account) List(db *gorm.DB, filter entity.Account) ([]entity.Account, error) {
	var accounts []entity.Account
	query := db.Model(&entity.Account{})

	if filter.UserID != uuid.Nil {
		query = query.Where("user_id = ?", filter.UserID)
	}
	if filter.Name != "" {
		query = query.Where("name = ?", filter.Name)
	}

	if err := query.Find(&accounts).Error; err != nil {
		return nil, err
	}

	return accounts, nil
}

func (a *Account) Save(db *gorm.DB, req entity.Account) error {
	return db.Save(req).Error
}
