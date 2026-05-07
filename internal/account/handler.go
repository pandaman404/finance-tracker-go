package account

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
	r.POST("/accounts/:userID", h.createAccount)
	r.GET("/accounts/:userID", h.getAccounts)
	r.PUT("/accounts/balance/:accountID", h.updateBalance)
	r.DELETE("/accounts/:accountID", h.deleteAccount)
}

func (h *Handler) createAccount(c *gin.Context) {
	var req CreateAccountRequest
	userID, err := uuid.Parse(c.Param("userID"))

	if err != nil {
		c.JSON(400, gin.H{"error": shared.ErrInvalidID.Error()})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": shared.ErrJSONBinding.Error()})
		return
	}

	account, err := h.service.CreateAccount(userID, req)

	if err != nil {
		switch err {
		case ErrInvalidName,
			ErrInvalidType:
			c.JSON(400, gin.H{"error": err.Error()})

		case ErrUserNotFound:
			c.JSON(404, gin.H{"error": err.Error()})

		case ErrAccountExists:
			c.JSON(409, gin.H{"error": err.Error()})

		default:
			c.JSON(500, gin.H{"error": shared.ErrInternalServer.Error()})
		}
		return
	}

	c.JSON(201, account)
}

func (h *Handler) getAccounts(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userID"))

	if err != nil {
		c.JSON(400, gin.H{"error": shared.ErrInvalidID.Error()})
		return
	}

	accounts, err := h.service.GetAccountsByUserID(userID)

	if err != nil {
		c.JSON(500, gin.H{"error": shared.ErrInternalServer.Error()})
		return
	}

	if len(accounts) == 0 {
		c.JSON(404, gin.H{"error": ErrAccountsNotFound.Error()})
		return
	}

	c.JSON(200, accounts)
}

func (h *Handler) updateBalance(c *gin.Context) {
	var req UpdateAccountBalanceRequest
	accountID, err := uuid.Parse(c.Param("accountID"))

	if err != nil {
		c.JSON(400, gin.H{"error": shared.ErrInvalidID.Error()})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": shared.ErrJSONBinding.Error()})
		return
	}

	account, err := h.service.UpdateAccountBalance(accountID, req)

	if err != nil {
		switch err {
		case ErrNotFound:
			c.JSON(404, gin.H{"error": err.Error()})
		default:
			c.JSON(500, gin.H{"error": shared.ErrInternalServer.Error()})
		}
	}

	c.JSON(200, account)
}

func (h *Handler) deleteAccount(c *gin.Context) {
	accountID, err := uuid.Parse(c.Param("accountID"))

	if err != nil {
		c.JSON(400, gin.H{"error": shared.ErrInvalidID.Error()})
		return
	}

	err = h.service.DeleteAccount(accountID)

	if err != nil {
		switch err {
		case ErrNotFound:
			c.JSON(404, gin.H{"error": err.Error()})
		default:
			c.JSON(500, gin.H{"error": shared.ErrInternalServer.Error()})
		}
		return
	}

	c.JSON(200, gin.H{"message": MsgDeleted})
}
