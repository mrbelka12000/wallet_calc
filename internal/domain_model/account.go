package domainmodel

import "github.com/google/uuid"

type (
	Account struct {
		ID      uuid.UUID
		UserID  uuid.UUID `validate:"required,uuid"`
		Name    string    `validate:"required"`
		Balance float64
	}
)
