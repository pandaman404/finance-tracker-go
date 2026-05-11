package transaction

import (
	"errors"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) *PostgresRepository {
	db.AutoMigrate(&Transaction{})
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(transaction *Transaction) error {
	return r.db.Create(transaction).Error
}

func (r *PostgresRepository) FindAll(userID uuid.UUID) ([]*Transaction, error) {
	var transactions []*Transaction

	result := r.db.
		Where("user_id = ?", userID).
		Order("transaction_date DESC").
		Find(&transactions)

	if result.Error != nil {
		return nil, result.Error
	}

	return transactions, nil
}

func (r *PostgresRepository) FindByID(id uuid.UUID, userID uuid.UUID) (*Transaction, error) {
	var transaction Transaction
	result := r.db.First(&transaction, "id = ? AND user_id = ?", id, userID)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &transaction, result.Error
}

func (r *PostgresRepository) Update(transaction *Transaction) error {
	return r.db.Save(transaction).Error
}

func (r *PostgresRepository) Delete(id uuid.UUID, userID uuid.UUID) error {
	result := r.db.Delete(&Transaction{}, "id = ? AND user_id = ?", id, userID)

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return result.Error
}

func (r *PostgresRepository) GetSummary(userID uuid.UUID, accountID *uuid.UUID) (*Summary, error) {
	type summaryResult struct {
		TotalIncome  decimal.Decimal
		TotalExpense decimal.Decimal
	}

	var result summaryResult

	query := r.db.Model(&Transaction{}).Where("user_id = ?", userID)

	if accountID != nil {
		query = query.Where("account_id = ?", *accountID)
	}

	err := query.
		Select(`
			COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) AS total_income,
			COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) AS total_expense
		`).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	balance := result.TotalIncome.Sub(result.TotalExpense)

	return &Summary{
		TotalIncome:  result.TotalIncome,
		TotalExpense: result.TotalExpense,
		Balance:      balance,
	}, nil
}

func (r *PostgresRepository) GetExpensesByCategory(userID uuid.UUID, accountID *uuid.UUID) ([]*ExpenseByCategory, error) {
	var expenses []*ExpenseByCategory

	query := r.db.Model(&Transaction{}).Where("type = ? AND user_id = ?", Expense, userID)

	if accountID != nil {
		query = query.Where("account_id = ?", *accountID)
	}

	err := query.
		Select("category_id, COALESCE(SUM(amount), 0) AS total").
		Group("category_id").
		Order("total DESC").
		Scan(&expenses).Error

	if err != nil {
		return nil, err
	}

	return expenses, nil
}
