package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/zalfrie/chatbot-ai/backend/config"
	"github.com/zalfrie/chatbot-ai/backend/routes"
	"log"
)

func main() {
	// Load .env file if present
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	e := echo.New()
	e.Use(middleware.Logger(), middleware.Recover())

	cfg := config.LoadConfig()
	db := config.InitDB(cfg)

	routes.Setup(e, db, cfg)

	e.Logger.Fatal(e.Start(cfg.ServerAddress))
}
