package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/pandaman404/finance-tracker-go/internal/account"
	"github.com/pandaman404/finance-tracker-go/internal/category"
	"github.com/pandaman404/finance-tracker-go/internal/config"
	"github.com/pandaman404/finance-tracker-go/internal/database"
	"github.com/pandaman404/finance-tracker-go/internal/middleware"
	"github.com/pandaman404/finance-tracker-go/internal/transaction"
	"github.com/pandaman404/finance-tracker-go/internal/user"
	"github.com/pandaman404/finance-tracker-go/pkg/logger"
	"golang.org/x/time/rate"
)

func main() {
	log := logger.New()

	if err := godotenv.Load("config/.env"); err != nil {
		log.Warn("no .env file found, using environment variables")
	}

	cfg := config.Load()
	db := database.NewPostgresDB(cfg)

	r := gin.New()
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS(cfg.CORSAllowedOrigins))
	r.Use(middleware.RateLimiter(rate.Limit(10), 20))

	// USER
	userRepo := user.NewPostgresRepository(db)
	userService := user.NewService(userRepo, cfg.JWTSecret)
	userHandler := user.NewHandler(userService)

	// ACCOUNT
	accountRepo := account.NewPostgresRepository(db)
	accountService := account.NewService(accountRepo, userRepo)
	accountHandler := account.NewHandler(accountService)

	// CATEGORY
	categoryRepo := category.NewPostgresRepository(db)
	categoryService := category.NewService(categoryRepo, userRepo)
	categoryHandler := category.NewHandler(categoryService)

	// TRANSACTION
	transactionRepo := transaction.NewPostgresRepository(db)
	transactionService := transaction.NewService(transactionRepo, accountRepo, categoryRepo)
	transactionHandler := transaction.NewHandler(transactionService)

	// Public routes
	userHandler.RegisterRoutes(r)

	// Login route with strict rate limit (2 req/s, burst 5)
	loginGroup := r.Group("/")
	loginGroup.Use(middleware.RateLimiter(rate.Limit(2), 5))
	{
		userHandler.RegisterLoginRoute(loginGroup)
	}

	// Protected routes
	protected := r.Group("/")
	protected.Use(middleware.Auth(cfg.JWTSecret))
	{
		accountHandler.RegisterRoutes(protected)
		categoryHandler.RegisterRoutes(protected)
		transactionHandler.RegisterRoutes(protected)
	}

	r.Run(":" + cfg.ServerPort)
}
