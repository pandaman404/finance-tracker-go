package transaction

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type TransactionType string

const (
	Income  TransactionType = "income"
	Expense TransactionType = "expense"
)

type Transaction struct {
	ID              uuid.UUID       `gorm:"type:uuid;primaryKey"`
	UserID          uuid.UUID       `gorm:"type:uuid;not null; index"`
	AccountID       uuid.UUID       `gorm:"type:uuid;not null; index"`
	CategoryID      uuid.UUID       `gorm:"type:uuid;not null; index"`
	Amount          decimal.Decimal `gorm:"not null;type:decimal(15,2)"`
	Type            TransactionType `gorm:"type:varchar(20);not null"`
	Description     string          `gorm:"type:varchar(255)"`
	TransactionDate time.Time       `gorm:"not null"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
