package transaction

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pandaman404/finance-tracker-go/internal/middleware"
	"github.com/pandaman404/finance-tracker-go/internal/shared"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) RegisterRoutes(r gin.IRouter) {
	r.POST("/transactions", h.createTransaction)
	r.GET("/transactions/summary", h.getSummary)
	r.GET("/transactions/expenses-by-category", h.getExpensesByCategory)
	r.GET("/transactions", h.getTransactions)
	r.GET("/transactions/:id", h.getTransactionByID)
	r.PUT("/transactions/:id", h.updateTransaction)
	r.DELETE("/transactions/:id", h.deleteTransaction)
}

func (h *Handler) createTransaction(c *gin.Context) {
	var req CreateTransactionRequest

	userID, ok := getCurrentUserID(c)
	if !ok {
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"errors": shared.ParseValidationErrors(err)})
		return
	}

	transaction, err := h.service.CreateTransaction(userID, req)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(201, transaction)
}

func (h *Handler) getTransactions(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		return
	}

	transactions, err := h.service.GetTransactions(userID)
	if err != nil {
		c.JSON(500, gin.H{"error": shared.ErrInternalServer.Error()})
		return
	}

	if len(transactions) == 0 {
		c.JSON(404, gin.H{"error": ErrTransactionsNotFound.Error()})
		return
	}

	c.JSON(200, transactions)
}

func (h *Handler) getTransactionByID(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": shared.ErrInvalidID.Error()})
		return
	}

	transaction, err := h.service.GetTransactionByID(id, userID)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(200, transaction)
}

func (h *Handler) updateTransaction(c *gin.Context) {
	var req UpdateTransactionRequest

	userID, ok := getCurrentUserID(c)
	if !ok {
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": shared.ErrInvalidID.Error()})
		return
	}

	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"errors": shared.ParseValidationErrors(err)})
		return
	}

	transaction, err := h.service.UpdateTransaction(id, userID, req)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(200, transaction)
}

func (h *Handler) deleteTransaction(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": shared.ErrInvalidID.Error()})
		return
	}

	err = h.service.DeleteTransaction(id, userID)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(200, gin.H{"message": MsgDeleted})
}

func (h *Handler) getSummary(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		return
	}

	accountID, ok := parseOptionalAccountID(c)
	if !ok {
		return
	}

	summary, err := h.service.GetSummary(userID, accountID)
	if err != nil {
		if err == ErrAccountNotFound {
			c.JSON(404, gin.H{"error": err.Error()})
			return
		}
		c.JSON(500, gin.H{"error": shared.ErrInternalServer.Error()})
		return
	}

	c.JSON(200, summary)
}

func (h *Handler) getExpensesByCategory(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		return
	}

	accountID, ok := parseOptionalAccountID(c)
	if !ok {
		return
	}

	expenses, err := h.service.GetExpensesByCategory(userID, accountID)
	if err != nil {
		if err == ErrAccountNotFound {
			c.JSON(404, gin.H{"error": err.Error()})
			return
		}
		c.JSON(500, gin.H{"error": shared.ErrInternalServer.Error()})
		return
	}

	c.JSON(200, expenses)
}

func parseOptionalAccountID(c *gin.Context) (*uuid.UUID, bool) {
	raw := c.Query("account_id")
	if raw == "" {
		return nil, true
	}

	accountID, err := uuid.Parse(raw)
	if err != nil {
		c.JSON(400, gin.H{"error": shared.ErrInvalidID.Error()})
		return nil, false
	}

	return &accountID, true
}

func getCurrentUserID(c *gin.Context) (uuid.UUID, bool) {
	raw, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(401, gin.H{"error": shared.ErrUnauthorized.Error()})
		return uuid.Nil, false
	}

	userID, err := uuid.Parse(raw.(string))
	if err != nil {
		c.JSON(400, gin.H{"error": shared.ErrInvalidID.Error()})
		return uuid.Nil, false
	}

	return userID, true
}

func (h *Handler) handleServiceError(c *gin.Context, err error) {
	switch err {
	case ErrInvalidAmount,
		ErrInvalidType,
		ErrInvalidDescription,
		ErrCategoryTypeMismatch,
		ErrCategoryNotFound,
		ErrAccountNotFound:
		c.JSON(400, gin.H{"error": err.Error()})
	case ErrTransactionNotFound:
		c.JSON(404, gin.H{"error": err.Error()})
	default:
		c.JSON(500, gin.H{"error": shared.ErrInternalServer.Error()})
	}
}
