package transaction

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreateTransactionRequest struct {
	AccountID       string          `json:"account_id" binding:"required"`
	CategoryID      string          `json:"category_id" binding:"required"`
	Amount          decimal.Decimal `json:"amount" binding:"required"`
	Type            TransactionType `json:"type" binding:"required"`
	Description     string          `json:"description" binding:"omitempty,max=255"`
	TransactionDate *time.Time      `json:"transaction_date" binding:"omitempty"`
}

type UpdateTransactionRequest struct {
	AccountID       *string          `json:"account_id" binding:"omitempty"`
	CategoryID      *string          `json:"category_id" binding:"omitempty"`
	Amount          *decimal.Decimal `json:"amount" binding:"omitempty"`
	Type            *TransactionType `json:"type" binding:"omitempty"`
	Description     *string          `json:"description" binding:"omitempty,max=255"`
	TransactionDate *time.Time       `json:"transaction_date" binding:"omitempty"`
}

type TransactionResponse struct {
	ID              string          `json:"id"`
	UserID          string          `json:"user_id"`
	AccountID       string          `json:"account_id"`
	CategoryID      string          `json:"category_id"`
	Amount          decimal.Decimal `json:"amount"`
	Type            TransactionType `json:"type"`
	Description     string          `json:"description"`
	TransactionDate string          `json:"transaction_date"`
}

type Summary struct {
	TotalIncome  decimal.Decimal `json:"total_income"`
	TotalExpense decimal.Decimal `json:"total_expense"`
	Balance      decimal.Decimal `json:"balance"`
}

type ExpenseByCategory struct {
	CategoryID uuid.UUID       `json:"category_id"`
	Total      decimal.Decimal `json:"total"`
}
