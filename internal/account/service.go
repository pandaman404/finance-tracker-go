package account

import (
	"github.com/google/uuid"
	"github.com/pandaman404/finance-tracker-go/internal/user"
)

type Service struct {
	accountRepo Repository
	userRepo    user.Repository
}

func NewService(accountRepo Repository, userRepo user.Repository) *Service {
	return &Service{accountRepo: accountRepo, userRepo: userRepo}
}

func (s *Service) CreateAccount(userID uuid.UUID, req CreateAccountRequest) (*AccountResponse, error) {

	switch req.Type {
	case Cash, Bank, CreditCard:
		// válido
	default:
		return nil, ErrInvalidType
	}

	user, err := s.userRepo.FindByID(userID)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	accounts, err := s.accountRepo.FindByUserID(userID)

	if err != nil {
		return nil, err
	}

	for _, acc := range accounts {
		if acc.Name == req.Name {
			return nil, ErrAccountExists
		}
	}

	account := &Account{
		ID:      uuid.New(),
		UserID:  userID,
		Name:    req.Name,
		Type:    req.Type,
		Balance: req.Balance,
	}

	if err = s.accountRepo.Create(account); err != nil {
		return nil, err
	}

	return toResponse(account), nil
}

func (s *Service) UpdateAccountBalance(accountID uuid.UUID, req UpdateAccountBalanceRequest) (*AccountResponse, error) {
	account, err := s.accountRepo.FindByID(accountID)

	if err != nil {
		return nil, err
	}

	if account == nil {
		return nil, ErrNotFound
	}

	updatedAccount, err := s.accountRepo.UpdateBalance(accountID, req.Amount)

	if err != nil {
		return nil, err
	}

	return toResponse(updatedAccount), nil
}

func (s *Service) GetAccountsByUserID(userID uuid.UUID) ([]*AccountResponse, error) {
	accounts, err := s.accountRepo.FindByUserID(userID)

	if err != nil {
		return nil, err
	}

	responses := make([]*AccountResponse, len(accounts))

	for i, acc := range accounts {
		responses[i] = toResponse(acc)
	}

	return responses, nil
}

func (s *Service) DeleteAccount(accountID uuid.UUID) error {
	account, err := s.accountRepo.FindByID(accountID)

	if err != nil {
		return err
	}

	if account == nil {
		return ErrNotFound
	}

	if err = s.accountRepo.Delete(accountID); err != nil {
		return err
	}

	return nil
}

func toResponse(a *Account) *AccountResponse {
	return &AccountResponse{
		ID:      a.ID.String(),
		UserID:  a.UserID.String(),
		Name:    a.Name,
		Type:    string(a.Type),
		Balance: a.Balance,
	}
}
