// Package dto (Data Transfer Object) - used for API/controller layer
// Example: UserDTO
package dto

import (
	"time"

	"github.com/google/uuid"
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
