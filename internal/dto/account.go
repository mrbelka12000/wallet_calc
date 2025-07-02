package dto

import (
	"time"

	"github.com/google/uuid"
)

type (
	AccountReq struct {
		ID   uuid.UUID `json:"id"`
		Name string    `json:"name"`
	}

	AccountResp struct {
		ID        uuid.UUID `json:"id"`
		UserID    uuid.UUID `json:"user_id"`
		Name      string    `json:"name"`
		Balance   float64   `json:"balance"`
		CreatedAt time.Time `json:"created_at"`
	}
)
