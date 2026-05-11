package transaction

import (
	"time"

	"github.com/google/uuid"
	"github.com/pandaman404/finance-tracker-go/internal/account"
	"github.com/pandaman404/finance-tracker-go/internal/category"
	"github.com/shopspring/decimal"
)

type Service struct {
	transactionRepo Repository
	accountRepo     account.Repository
	categoryRepo    category.Repository
}

func NewService(
	transactionRepo Repository,
	accountRepo account.Repository,
	categoryRepo category.Repository,
) *Service {
	return &Service{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
		categoryRepo:    categoryRepo,
	}
}

func (s *Service) CreateTransaction(userID uuid.UUID, req CreateTransactionRequest) (*TransactionResponse, error) {
	if req.Amount.LessThanOrEqual(decimal.Zero) {
		return nil, ErrInvalidAmount
	}

	accountID, err := uuid.Parse(req.AccountID)
	if err != nil {
		return nil, ErrAccountNotFound
	}

	categoryID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		return nil, ErrCategoryNotFound
	}

	if err = s.validateOwnershipAndType(userID, accountID, categoryID, req.Type); err != nil {
		return nil, err
	}

	transactionDate := time.Now()
	if req.TransactionDate != nil {
		transactionDate = *req.TransactionDate
	}

	transaction := &Transaction{
		ID:              uuid.New(),
		UserID:          userID,
		AccountID:       accountID,
		CategoryID:      categoryID,
		Amount:          req.Amount,
		Type:            req.Type,
		Description:     req.Description,
		TransactionDate: transactionDate,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err = s.transactionRepo.Create(transaction); err != nil {
		return nil, err
	}

	return toResponse(transaction), nil
}

func (s *Service) GetTransactions(userID uuid.UUID) ([]*TransactionResponse, error) {
	transactions, err := s.transactionRepo.FindAll(userID)
	if err != nil {
		return nil, err
	}

	responses := make([]*TransactionResponse, len(transactions))
	for i, t := range transactions {
		responses[i] = toResponse(t)
	}

	return responses, nil
}

func (s *Service) GetTransactionByID(id uuid.UUID, userID uuid.UUID) (*TransactionResponse, error) {
	transaction, err := s.transactionRepo.FindByID(id, userID)
	if err != nil {
		return nil, err
	}

	if transaction == nil {
		return nil, ErrTransactionNotFound
	}

	return toResponse(transaction), nil
}

func (s *Service) UpdateTransaction(id uuid.UUID, userID uuid.UUID, req UpdateTransactionRequest) (*TransactionResponse, error) {
	transaction, err := s.transactionRepo.FindByID(id, userID)
	if err != nil {
		return nil, err
	}

	if transaction == nil {
		return nil, ErrTransactionNotFound
	}

	if req.AccountID != nil {
		accountID, parseErr := uuid.Parse(*req.AccountID)
		if parseErr != nil {
			return nil, ErrAccountNotFound
		}
		transaction.AccountID = accountID
	}

	if req.CategoryID != nil {
		categoryID, parseErr := uuid.Parse(*req.CategoryID)
		if parseErr != nil {
			return nil, ErrCategoryNotFound
		}
		transaction.CategoryID = categoryID
	}

	if req.Amount != nil {
		if req.Amount.LessThanOrEqual(decimal.Zero) {
			return nil, ErrInvalidAmount
		}
		transaction.Amount = *req.Amount
	}

	if req.Type != nil {
		transaction.Type = *req.Type
	}

	if req.Description != nil {
		transaction.Description = *req.Description
	}

	if req.TransactionDate != nil {
		transaction.TransactionDate = *req.TransactionDate
	}

	if err = s.validateOwnershipAndType(userID, transaction.AccountID, transaction.CategoryID, transaction.Type); err != nil {
		return nil, err
	}

	transaction.UpdatedAt = time.Now()

	if err = s.transactionRepo.Update(transaction); err != nil {
		return nil, err
	}

	return toResponse(transaction), nil
}

func (s *Service) DeleteTransaction(id uuid.UUID, userID uuid.UUID) error {
	transaction, err := s.transactionRepo.FindByID(id, userID)
	if err != nil {
		return err
	}

	if transaction == nil {
		return ErrTransactionNotFound
	}

	if err = s.transactionRepo.Delete(id, userID); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetSummary(userID uuid.UUID, accountID *uuid.UUID) (*Summary, error) {
	if accountID != nil {
		acc, err := s.accountRepo.FindByID(*accountID)
		if err != nil {
			return nil, err
		}
		if acc == nil || acc.UserID != userID {
			return nil, ErrAccountNotFound
		}
	}

	return s.transactionRepo.GetSummary(userID, accountID)
}

func (s *Service) GetExpensesByCategory(userID uuid.UUID, accountID *uuid.UUID) ([]*ExpenseByCategory, error) {
	if accountID != nil {
		acc, err := s.accountRepo.FindByID(*accountID)
		if err != nil {
			return nil, err
		}
		if acc == nil || acc.UserID != userID {
			return nil, ErrAccountNotFound
		}
	}

	return s.transactionRepo.GetExpensesByCategory(userID, accountID)
}

func (s *Service) validateOwnershipAndType(
	userID uuid.UUID,
	accountID uuid.UUID,
	categoryID uuid.UUID,
	typeTx TransactionType,
) error {
	accountEntity, err := s.accountRepo.FindByID(accountID)
	if err != nil {
		return err
	}

	if accountEntity == nil || accountEntity.UserID != userID {
		return ErrAccountNotFound
	}

	categoryEntity, err := s.categoryRepo.FindByID(categoryID)
	if err != nil {
		return err
	}

	if categoryEntity == nil {
		return ErrCategoryNotFound
	}

	if categoryEntity.UserID != nil && *categoryEntity.UserID != userID {
		return ErrCategoryNotFound
	}

	if category.CategoryType(typeTx) != categoryEntity.Type {
		return ErrCategoryTypeMismatch
	}

	return nil
}

func toResponse(t *Transaction) *TransactionResponse {
	return &TransactionResponse{
		ID:              t.ID.String(),
		UserID:          t.UserID.String(),
		AccountID:       t.AccountID.String(),
		CategoryID:      t.CategoryID.String(),
		Amount:          t.Amount,
		Type:            t.Type,
		Description:     t.Description,
		TransactionDate: t.TransactionDate.Format(time.RFC3339),
	}
}
