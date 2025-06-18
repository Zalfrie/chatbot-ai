package routes

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/zalfrie/chatbot-ai/backend/config"
	"github.com/zalfrie/chatbot-ai/backend/controllers"
	"github.com/zalfrie/chatbot-ai/backend/middleware"
)

func Setup(e *echo.Echo, db *sqlx.DB, cfg *config.Config) {
	// Public routes
	e.POST("/api/register", controllers.Register(db))
	e.POST("/api/login", controllers.Login(db, cfg.JWTSecret))
	e.POST("/api/forgot-password", controllers.ForgotPassword(db, cfg))
	e.POST("/api/reset-password", controllers.ResetPassword(db))

	// Protected
	jwt := middleware.JWTMiddleware(cfg.JWTSecret)
	g := e.Group("/api", jwt)

	g.GET("/memory", controllers.MemoryList(db))                                // all authenticated users
	g.DELETE("/memory/:id", controllers.DeleteMemory(db), middleware.AdminOnly) // only admin
	g.GET("/ws", controllers.WebSocketHandler(db))
}
