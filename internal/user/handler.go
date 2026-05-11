package user

import (
	"strings"

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
	r.POST("/users", h.createUser)
	r.GET("/users", h.getUsers)
	r.GET("/users/:id", h.getUserByID)
	r.PUT("/users/:id", h.updateUser)
	r.DELETE("/users/:id", h.deleteUser)
}

func (h *Handler) RegisterLoginRoute(r gin.IRouter) {
	r.POST("/users/login", h.login)
}

func (h *Handler) login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"errors": shared.ParseValidationErrors(err)})
		return
	}

	token, err := h.service.Login(req)
	if err != nil {
		switch err {
		case ErrInvalidCredentials:
			c.JSON(401, gin.H{"error": err.Error()})
		default:
			c.JSON(500, gin.H{"error": shared.ErrInternalServer.Error()})
		}
		return
	}

	c.JSON(200, token)
}

func (h *Handler) createUser(c *gin.Context) {
	var req CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"errors": shared.ParseValidationErrors(err)})
		return
	}

	user, err := h.service.CreateUser(req)
	if err != nil {
		if strings.Contains(err.Error(), ErrEmailExists.Error()) {
			c.JSON(409, gin.H{"error": err.Error()})
			return
		}
		c.JSON(500, gin.H{"error": shared.ErrInternalServer.Error()})
		return
	}

	c.JSON(201, user)
}

func (h *Handler) getUsers(c *gin.Context) {
	email := c.Query("email")

	if email != "" {
		user, err := h.service.GetUserByEmail(email)
		if err != nil {
			c.JSON(500, gin.H{"error": shared.ErrInternalServer.Error()})
			return
		}

		if user == nil {
			c.JSON(404, gin.H{"error": ErrNotFound.Error()})
			return
		}

		c.JSON(200, user)
		return
	}

	users, err := h.service.GetUsers()
	if err != nil {
		c.JSON(500, gin.H{"error": shared.ErrInternalServer.Error()})
		return
	}

	c.JSON(200, users)
}

func (h *Handler) getUserByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(400, gin.H{"error": shared.ErrInvalidID.Error()})
		return
	}

	user, err := h.service.GetUserByID(id)
	if err != nil {
		c.JSON(500, gin.H{"error": shared.ErrInternalServer.Error()})
		return
	}
	if user == nil {
		c.JSON(404, gin.H{"error": ErrNotFound.Error()})
		return
	}

	c.JSON(200, user)
}

func (h *Handler) updateUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(400, gin.H{"error": shared.ErrInvalidID.Error()})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"errors": shared.ParseValidationErrors(err)})
		return
	}

	user, err := h.service.UpdateUser(id, req)
	if err != nil {
		if strings.Contains(err.Error(), ErrEmailExists.Error()) {
			c.JSON(409, gin.H{"error": err.Error()})
			return
		}
		c.JSON(500, gin.H{"error": shared.ErrInternalServer.Error()})
		return
	}
	if user == nil {
		c.JSON(404, gin.H{"error": ErrNotFound.Error()})
		return
	}

	c.JSON(200, user)
}

func (h *Handler) deleteUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(400, gin.H{"error": shared.ErrInvalidID.Error()})
		return
	}

	deleted, err := h.service.DeleteUser(id)
	if err != nil {
		c.JSON(500, gin.H{"error": shared.ErrInternalServer.Error()})
		return
	}
	if !deleted {
		c.JSON(404, gin.H{"error": ErrNotFound.Error()})
		return
	}

	c.JSON(204, nil)
}
