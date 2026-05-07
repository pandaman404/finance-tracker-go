package account

import (
	"time"

	"github.com/google/uuid"
	"github.com/pandaman404/finance-tracker-go/internal/user"
	"github.com/shopspring/decimal"
)

type AccountType string

const (
	Cash       AccountType = "cash"
	Bank       AccountType = "bank"
	CreditCard AccountType = "credit_card"
)

type Account struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID uuid.UUID `gorm:"type:uuid;not null; index"`
	User   user.User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`

	Name      string          `gorm:"not null"`
	Type      AccountType     `gorm:"type:varchar(20);not null"`
	Balance   decimal.Decimal `gorm:"not null;type:decimal(15,2);default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
