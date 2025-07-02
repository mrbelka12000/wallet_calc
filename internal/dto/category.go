package dto

import "github.com/google/uuid"

type (
	CategoryReq struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}

	CategoryResp struct {
		ID   uuid.UUID `json:"id"`
		Name string    `json:"name"`
		Type string    `json:"type"`
	}
)
