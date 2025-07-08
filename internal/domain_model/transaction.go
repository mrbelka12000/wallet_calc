package domainmodel

import (
	"time"

	"github.com/google/uuid"
)

type (
	Transaction struct {
		ID          uuid.UUID
		UserID      uuid.UUID
		Date        string
		Description string
		Amount      float64
		CategoryID  uuid.UUID
		AccountID   uuid.UUID
		Type        string
		CreatedAt   time.Time
	}

	BaseTransaction struct {
		Date        string
		Description string
		Amount      float64
		Transaction string
		Category    string
		Confidence  float64
	}
)
