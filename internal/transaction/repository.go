package transaction

import (
	"github.com/google/uuid"
)

type Repository interface {
	Create(transaction *Transaction) error
	FindAll(userID uuid.UUID) ([]*Transaction, error)
	FindByID(id uuid.UUID, userID uuid.UUID) (*Transaction, error)
	Update(transaction *Transaction) error
	Delete(id uuid.UUID, userID uuid.UUID) error

	GetSummary(userID uuid.UUID, accountID *uuid.UUID) (*Summary, error)
	GetExpensesByCategory(userID uuid.UUID, accountID *uuid.UUID) ([]*ExpenseByCategory, error)
}
