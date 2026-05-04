package user

import (
	"errors"

	"gorm.io/gorm"
)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) *PostgresRepository {
	db.AutoMigrate(&User{})
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(user *User) error {
	return r.db.Create(user).Error
}

func (r *PostgresRepository) FindAll() ([]User, error) {
	var users []User
	result := r.db.Find(&users)
	return users, result.Error
}

func (r *PostgresRepository) FindByID(id string) (*User, error) {
	var user User
	result := r.db.First(&user, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, result.Error
}

func (r *PostgresRepository) FindByEmail(email string) (*User, error) {
	var user User
	result := r.db.Where("email = ?", email).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, result.Error
}

func (r *PostgresRepository) Update(user *User) error {
	return r.db.Save(user).Error
}

func (r *PostgresRepository) Delete(id string) error {
	return r.db.Delete(&User{}, "id = ?", id).Error
}
