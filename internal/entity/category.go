package entity

import "github.com/google/uuid"

type (
	Category struct {
		ID   uuid.UUID `gorm:"column:id;primary_key"`
		Name string    `gorm:"column:name"`
		Type string    `gorm:"column:type"`
	}
)

func (Category) TableName() string {
	return "categories"
}
