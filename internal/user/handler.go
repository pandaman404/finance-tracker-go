package user

import (
	"strings"

	"github.com/gin-gonic/gin"
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

func (h *Handler) createUser(c *gin.Context) {
	var req CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.CreateUser(req)
	if err != nil {
		if strings.Contains(err.Error(), "el email ya está registrado") {
			c.JSON(409, gin.H{"error": err.Error()})
			return
		}
		c.JSON(500, gin.H{"error": "error interno del servidor"})
		return
	}

	c.JSON(201, user)
}

func (h *Handler) getUsers(c *gin.Context) {
	users, err := h.service.GetUsers()
	if err != nil {
		c.JSON(500, gin.H{"error": "error interno del servidor"})
		return
	}

	c.JSON(200, users)
}

func (h *Handler) getUserByID(c *gin.Context) {
	id := c.Param("id")

	user, err := h.service.GetUserByID(id)
	if err != nil {
		c.JSON(500, gin.H{"error": "error interno del servidor"})
		return
	}
	if user == nil {
		c.JSON(404, gin.H{"error": "usuario no encontrado"})
		return
	}

	c.JSON(200, user)
}

func (h *Handler) updateUser(c *gin.Context) {
	id := c.Param("id")

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.UpdateUser(id, req)
	if err != nil {
		if strings.Contains(err.Error(), "el email ya está registrado") {
			c.JSON(409, gin.H{"error": err.Error()})
			return
		}
		c.JSON(500, gin.H{"error": "error interno del servidor"})
		return
	}
	if user == nil {
		c.JSON(404, gin.H{"error": "usuario no encontrado"})
		return
	}

	c.JSON(200, user)
}

func (h *Handler) deleteUser(c *gin.Context) {
	id := c.Param("id")

	deleted, err := h.service.DeleteUser(id)
	if err != nil {
		c.JSON(500, gin.H{"error": "error interno del servidor"})
		return
	}
	if !deleted {
		c.JSON(404, gin.H{"error": "usuario no encontrado"})
		return
	}

	c.JSON(204, nil)
}
