package category

import "github.com/google/uuid"

type Repository interface {
	Create(category *Category) error
	FindByID(id uuid.UUID) (*Category, error)
	Update(category *Category) error
	Delete(id uuid.UUID) error

	FindByUserID(userID uuid.UUID) ([]*Category, error)
	FindAvailableByUserID(userID uuid.UUID) ([]*Category, error)
	FindByNameAndUserID(name string, userID *uuid.UUID) (*Category, error)
}
