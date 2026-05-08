package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/pandaman404/finance-tracker-go/internal/account"
	"github.com/pandaman404/finance-tracker-go/internal/category"
	"github.com/pandaman404/finance-tracker-go/internal/config"
	"github.com/pandaman404/finance-tracker-go/internal/database"
	"github.com/pandaman404/finance-tracker-go/internal/user"
)

func main() {
	if err := godotenv.Load("config/.env"); err != nil {
		log.Println("no .env file found, using environment variables")
	}

	cfg := config.Load()
	db := database.NewPostgresDB(cfg)

	r := gin.Default()

	// USER
	userRepo := user.NewPostgresRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	// ACCOUNT
	accountRepo := account.NewPostgresRepository(db)
	accountService := account.NewService(accountRepo, userRepo)
	accountHandler := account.NewHandler(accountService)

	// CATEGORY
	categoryRepo := category.NewPostgresRepository(db)
	categoryService := category.NewService(categoryRepo, userRepo)
	categoryHandler := category.NewHandler(categoryService)

	// Routes
	userHandler.RegisterRoutes(r)
	accountHandler.RegisterRoutes(r)
	categoryHandler.RegisterRoutes(r)

	r.Run(":" + cfg.ServerPort)
}
