package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

	repo := user.NewPostgresRepository(db)
	service := user.NewService(repo)
	handler := user.NewHandler(service)

	handler.RegisterRoutes(r)

	r.Run(":" + cfg.ServerPort)
}
