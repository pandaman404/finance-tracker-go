package account

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Repository interface {
	Create(account *Account) error
	FindByID(id uuid.UUID) (*Account, error)
	Update(account *Account) error
	Delete(id uuid.UUID) error

	FindByUserID(userID uuid.UUID) ([]*Account, error)
	UpdateBalance(accountID uuid.UUID, amount decimal.Decimal) (*Account, error)
}
