package account

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
	db.AutoMigrate(&Account{})
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(account *Account) error {
	return r.db.Create(account).Error
}

func (r *PostgresRepository) FindByID(id uuid.UUID) (*Account, error) {
	var account Account
	result := r.db.First(&account, "id = ?", id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &account, result.Error
}

func (r *PostgresRepository) Update(account *Account) error {
	return r.db.Save(account).Error
}

func (r *PostgresRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&Account{}, "id = ?", id).Error
}

func (r *PostgresRepository) FindByUserID(userID uuid.UUID) ([]*Account, error) {
	var accounts []*Account

	result := r.db.Where("user_id = ?", userID).Find(&accounts)

	if result.Error != nil {
		return nil, result.Error
	}

	return accounts, nil

}

func (r *PostgresRepository) UpdateBalance(accountID uuid.UUID, amount decimal.Decimal) (*Account, error) {
	var account Account
	result := r.db.Model(&Account{}).
		Where("id = ?", accountID).
		Update("balance", gorm.Expr("balance + ?", amount)).
		First(&account, "id = ?", accountID)

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &account, result.Error
}
