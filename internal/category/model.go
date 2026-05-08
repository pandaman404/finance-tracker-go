package category

import (
	"time"

	"github.com/google/uuid"
)

type CategoryType string

const (
	Income  CategoryType = "income"
	Expense CategoryType = "expense"
)

type Category struct {
	ID        uuid.UUID    `gorm:"type:uuid;primaryKey"`
	UserID    *uuid.UUID   `gorm:"type:uuid; index"`
	Name      string       `gorm:"not null"`
	Type      CategoryType `gorm:"type:varchar(20);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
