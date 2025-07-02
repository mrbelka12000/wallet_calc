// Package entity (DB Model) – user for repo layer
// Example: UserEntity
package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	domainmodel "github.com/mrbelka12000/wallet_calc/internal/domain_model"
)

type (
	User struct {
		ID        uuid.UUID `gorm:"column:id"`
		Email     string    `gorm:"column:email"`
		Password  string    `gorm:"column:password"`
		Name      string    `gorm:"column:name"`
		CreatedAt time.Time `gorm:"column:created_at"`
	}
)

func (u *User) TableName() string {
	return "users"
}

func ConvertFromDMUser(user domainmodel.User) User {
	return User{
		ID:        user.ID,
		Email:     user.Email,
		Password:  user.Password,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
	}
}

func (u User) Matches(args interface{}) bool {
	arg, ok := args.(User)
	if !ok {
		return false
	}

	return u.Email == arg.Email && u.Name == arg.Name
}

func (u User) String() string {
	return fmt.Sprintf("User{Email: %s, Name: %s}", u.Email, u.Name)
}
