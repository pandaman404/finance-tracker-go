package category

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pandaman404/finance-tracker-go/internal/shared"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.POST("/categories/:userID", h.createCategory)
	r.GET("/categories/:userID", h.getCategories)
	r.PUT("/categories/:categoryID", h.updateCategory)
	r.DELETE("/categories/:categoryID", h.deleteCategory)
}

func (h *Handler) createCategory(c *gin.Context) {
	var req CreateCategoryRequest
	userID, err := uuid.Parse(c.Param("userID"))

	if err != nil {
		c.JSON(400, gin.H{"error": shared.ErrInvalidID.Error()})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"errors": shared.ParseValidationErrors(err)})
		return
	}

	category, err := h.service.CreateCategory(userID, req)

	if err != nil {
		switch err {
		case ErrInvalidType:
			c.JSON(400, gin.H{"error": err.Error()})

		case ErrUserNotFound:
			c.JSON(404, gin.H{"error": err.Error()})

		case ErrCategoryExists:
			c.JSON(409, gin.H{"error": err.Error()})

		default:
			c.JSON(500, gin.H{"error": shared.ErrInternalServer.Error()})
		}
		return
	}

	c.JSON(201, category)
}

func (h *Handler) getCategories(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userID"))

	if err != nil {
		c.JSON(400, gin.H{"error": shared.ErrInvalidID.Error()})
		return
	}

	categories, err := h.service.GetAvailableCategories(userID)

	if err != nil {
		c.JSON(500, gin.H{"error": shared.ErrInternalServer.Error()})
		return
	}

	if len(categories) == 0 {
		c.JSON(404, gin.H{"error": ErrCategoriesNotFound.Error()})
		return
	}

	c.JSON(200, categories)
}

func (h *Handler) updateCategory(c *gin.Context) {
	var req UpdateCategoryRequest
	categoryID, err := uuid.Parse(c.Param("categoryID"))

	if err != nil {
		c.JSON(400, gin.H{"error": shared.ErrInvalidID.Error()})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"errors": shared.ParseValidationErrors(err)})
		return
	}

	category, err := h.service.UpdateCategory(categoryID, req)

	if err != nil {
		switch err {
		case ErrInvalidType:
			c.JSON(400, gin.H{"error": err.Error()})

		case ErrCategoryNotFound:
			c.JSON(404, gin.H{"error": err.Error()})

		default:
			c.JSON(500, gin.H{"error": shared.ErrInternalServer.Error()})
		}
		return
	}

	c.JSON(200, category)
}

func (h *Handler) deleteCategory(c *gin.Context) {
	categoryID, err := uuid.Parse(c.Param("categoryID"))

	if err != nil {
		c.JSON(400, gin.H{"error": shared.ErrInvalidID.Error()})
		return
	}

	err = h.service.DeleteCategory(categoryID)

	if err != nil {
		switch err {
		case ErrCategoryNotFound:
			c.JSON(404, gin.H{"error": err.Error()})
		default:
			c.JSON(500, gin.H{"error": shared.ErrInternalServer.Error()})
		}
		return
	}

	c.JSON(200, gin.H{"message": MsgDeleted})
}
