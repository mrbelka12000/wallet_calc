package repo

import (
	"gorm.io/gorm"

	"github.com/mrbelka12000/wallet_calc/internal/entity"
)

type (
	Category struct {
	}
)

func NewCategory() *Category {
	return &Category{}
}

func (c *Category) Create(db *gorm.DB, category entity.Category) error {
	return db.Create(&category).Error
}

func (c *Category) List(db *gorm.DB) ([]entity.Category, error) {
	var categories []entity.Category
	
	if err := db.Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}
