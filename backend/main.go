package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/zalfrie/chatbot-ai/backend/config"
	"github.com/zalfrie/chatbot-ai/backend/routes"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger(), middleware.Recover())

	cfg := config.LoadConfig()
	db := config.InitDB(cfg)

	routes.Setup(e, db, cfg)

	e.Logger.Fatal(e.Start(cfg.ServerAddress))
}
