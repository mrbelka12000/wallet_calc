// Package dto (Data Transfer Object) - used for API/controller layer
// Example: UserDTO
package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/mrbelka12000/wallet_calc/internal/entity"
)

type (
	UserCU struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}

	User struct {
		ID        uuid.UUID `json:"id,omitempty"`
		Email     string    `json:"email,omitempty"`
		Password  string    `json:"password,omitempty"`
		Name      string    `json:"name,omitempty"`
		CreatedAt time.Time `json:"created_at"`
	}

	UserLogin struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
)

func ConvertFromEntityUser(u entity.User, withPassword bool) User {
	user := User{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
	}

	if withPassword {
		user.Password = u.Password
	}

	return user
}
