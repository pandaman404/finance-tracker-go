package category

type CreateCategoryRequest struct {
	Name string       `json:"name" binding:"required,max=100"`
	Type CategoryType `json:"type" binding:"required"`
}

type UpdateCategoryRequest struct {
	Name string       `json:"name" binding:"omitempty,max=100"`
	Type CategoryType `json:"type" binding:"omitempty"`
}

type CategoryResponse struct {
	ID     string       `json:"id"`
	UserID *string      `json:"userID,omitempty"`
	Name   string       `json:"name"`
	Type   CategoryType `json:"type"`
}
