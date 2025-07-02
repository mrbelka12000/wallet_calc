package domainmodel

import "github.com/google/uuid"

type (
	Category struct {
		ID   uuid.UUID
		Name string `validate:"required"`
		Type string `validate:"required"`
	}
)
