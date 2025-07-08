package entity

import (
	"time"

	"github.com/google/uuid"
)

type (
	Transaction struct {
		ID          uuid.UUID `gorm:"column:id;primary_key"`
		UserID      uuid.UUID `gorm:"column:user_id"`
		Date        string    `gorm:"column:date"`
		Description string    `gorm:"column:description"`
		Amount      float64   `gorm:"column:amount"`
		CategoryID  uuid.UUID `gorm:"column:category_id"`
		AccountID   uuid.UUID `gorm:"column:account_id"`
		Type        string    `gorm:"column:type"`
		CreatedAt   time.Time `gorm:"column:created_at"`
	}
)

func (u *Transaction) TableName() string {
	return "transactions"
}
