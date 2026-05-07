package account

import "github.com/shopspring/decimal"

type CreateAccountRequest struct {
	Name    string          `json:"name" binding:"required,max=100"`
	Type    AccountType     `json:"type" binding:"required"`
	Balance decimal.Decimal `json:"balance" binding:"omitempty"`
}

type UpdateAccountBalanceRequest struct {
	Amount decimal.Decimal `json:"amount" binding:"required"`
}

type AccountResponse struct {
	ID      string          `json:"id"`
	UserID  string          `json:"user_id"`
	Name    string          `json:"name"`
	Type    string          `json:"type"`
	Balance decimal.Decimal `json:"balance"`
}
