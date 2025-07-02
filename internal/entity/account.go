package entity

import (
	"time"

	"github.com/google/uuid"
)

type (
	Account struct {
		ID        uuid.UUID `gorm:"column:id;primary_key"`
		UserID    uuid.UUID `gorm:"column:user_id"`
		Name      string    `gorm:"column:name"`
		Balance   float64   `gorm:"column:balance"`
		CreatedAt time.Time `gorm:"column:created_at"`
	}
)

func (Account) TableName() string {
	return "accounts"
}
