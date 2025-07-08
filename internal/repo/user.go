package repo

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/mrbelka12000/wallet_calc/internal/entity"
)

type User struct {
}

func NewUser() *User {
	return &User{}
}

func (u *User) Create(db *gorm.DB, req entity.User) error {
	return db.Create(&req).Error
}

func (u *User) Get(db *gorm.DB, req entity.User) (entity.User, error) {
	var result entity.User
	query := db.Model(&entity.User{})

	// Dynamically build the query from non-zero fields
	if req.ID != uuid.Nil {
		query = query.Where("id = ?", req.ID)
	}
	if req.Email != "" {
		query = query.Where("email = ?", req.Email)
	}
	if req.Password != "" {
		query = query.Where("password = ?", req.Password)
	}
	if req.Name != "" {
		query = query.Where("name = ?", req.Name)
	}

	if err := query.First(&result).Error; err != nil {
		return result, err
	}

	return result, nil
}
