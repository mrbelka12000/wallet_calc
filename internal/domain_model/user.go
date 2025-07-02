// Package domainmodel Business Logic Layer - used for usecase/service layer
// Example: User
package domainmodel

import (
	"time"

	"github.com/google/uuid"
)

type (
	User struct {
		ID        uuid.UUID
		Email     string `validate:"required,email"`
		Password  string `validate:"required"`
		Name      string `validate:"required"`
		CreatedAt time.Time
	}
)
