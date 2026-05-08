package category

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) *PostgresRepository {
	db.AutoMigrate(&Category{})
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(category *Category) error {
	return r.db.Create(category).Error
}

func (r *PostgresRepository) FindByID(id uuid.UUID) (*Category, error) {
	var category Category
	result := r.db.First(&category, "id = ?", id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &category, result.Error
}

func (r *PostgresRepository) Update(category *Category) error {
	return r.db.Save(category).Error
}

func (r *PostgresRepository) Delete(id uuid.UUID) error {
	result := r.db.Delete(&Category{}, "id = ?", id)

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return result.Error
}

func (r *PostgresRepository) FindByUserID(userID uuid.UUID) ([]*Category, error) {
	var categories []*Category

	result := r.db.Where("user_id = ?", userID).Find(&categories)

	if result.Error != nil {
		return nil, result.Error
	}

	return categories, nil
}

func (r *PostgresRepository) FindAvailableByUserID(userID uuid.UUID) ([]*Category, error) {
	var categories []*Category

	result := r.db.
		Where("user_id = ? OR user_id IS NULL", userID).
		Find(&categories)

	if result.Error != nil {
		return nil, result.Error
	}

	return categories, nil
}

func (r *PostgresRepository) FindByNameAndUserID(name string, userID *uuid.UUID) (*Category, error) {
	var category Category

	query := r.db.Where("name = ?", name)

	if userID == nil {
		query = query.Where("user_id IS NULL")
	} else {
		query = query.Where("user_id = ?", *userID)
	}

	result := query.First(&category)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &category, result.Error
}
