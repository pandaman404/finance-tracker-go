package category

import (
	"github.com/google/uuid"
	"github.com/pandaman404/finance-tracker-go/internal/user"
)

type Service struct {
	categoryRepo Repository
	userRepo     user.Repository
}

func NewService(categoryRepo Repository, userRepo user.Repository) *Service {
	return &Service{categoryRepo: categoryRepo, userRepo: userRepo}
}

func (s *Service) CreateCategory(userID uuid.UUID, req CreateCategoryRequest) (*CategoryResponse, error) {
	switch req.Type {

	case Income, Expense:
		// válido
	default:
		return nil, ErrInvalidType
	}

	categories, err := s.categoryRepo.FindByUserID(userID)

	if err != nil {
		return nil, err
	}

	for _, cat := range categories {
		if cat.Name == req.Name {
			return nil, ErrCategoryExists
		}
	}

	category := &Category{
		ID:     uuid.New(),
		UserID: &userID,
		Name:   req.Name,
		Type:   req.Type,
	}

	if err = s.categoryRepo.Create(category); err != nil {
		return nil, err
	}

	return toResponse(category), nil
}

func (s *Service) GetAvailableCategories(userID uuid.UUID) ([]*CategoryResponse, error) {
	categories, err := s.categoryRepo.FindAvailableByUserID(userID)

	if err != nil {
		return nil, err
	}

	responses := make([]*CategoryResponse, len(categories))

	for i, cat := range categories {
		responses[i] = toResponse(cat)
	}

	return responses, nil
}

func (s *Service) UpdateCategory(categoryID uuid.UUID, req UpdateCategoryRequest) (*CategoryResponse, error) {
	category, err := s.categoryRepo.FindByID(categoryID)

	if err != nil {
		return nil, err
	}

	if category == nil {
		return nil, ErrCategoryNotFound
	}

	if req.Name != "" {
		category.Name = req.Name
	}

	if req.Type != "" {
		switch req.Type {
		case Income, Expense:
			category.Type = req.Type
		default:
			return nil, ErrInvalidType
		}
	}

	existingCategory, err := s.categoryRepo.FindByNameAndUserID(
		category.Name,
		category.UserID,
	)

	if err != nil {
		return nil, err
	}

	if existingCategory != nil && existingCategory.ID != category.ID {
		return nil, ErrCategoryExists
	}

	if err = s.categoryRepo.Update(category); err != nil {
		return nil, err
	}

	return toResponse(category), nil
}

func (s *Service) DeleteCategory(categoryID uuid.UUID) error {
	category, err := s.categoryRepo.FindByID(categoryID)

	if err != nil {
		return err
	}

	if category == nil {
		return ErrCategoryNotFound
	}

	if err = s.categoryRepo.Delete(categoryID); err != nil {
		return err
	}

	return nil
}

func toResponse(cat *Category) *CategoryResponse {
	var userID *string

	if cat.UserID != nil {
		id := cat.UserID.String()
		userID = &id
	}

	return &CategoryResponse{
		ID:     cat.ID.String(),
		UserID: userID,
		Name:   cat.Name,
		Type:   cat.Type,
	}
}
