package controllers

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/zalfrie/chatbot-ai/backend/config"
	"github.com/zalfrie/chatbot-ai/backend/models"
	"github.com/zalfrie/chatbot-ai/backend/utils"
)

// Register new user (default role = user)
func Register(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to hash password")
		}
		// Insert user with default role "user"
		_, err = db.NamedExec(`
            INSERT INTO users (name, email, password, role, created_at)
            VALUES (:name, :email, :password, 'user', :created_at)
        `, map[string]interface{}{
			"name":       req.Name,
			"email":      req.Email,
			"password":   string(hashed),
			"created_at": time.Now(),
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusCreated, echo.Map{"message": "Registered successfully"})
	}
}

// Login issues JWT
func Login(db *sqlx.DB, cfg *config.Config) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		var user models.User
		if err := db.Get(&user, "SELECT * FROM users WHERE email=?", req.Email); err != nil {
			return echo.ErrUnauthorized
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			return echo.ErrUnauthorized
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":   user.ID,
			"role": user.Role,
			"exp":  time.Now().Add(72 * time.Hour).Unix(),
		})
		t, err := token.SignedString([]byte(cfg.JWTSecret))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to sign token")
		}
		return c.JSON(http.StatusOK, echo.Map{"token": t})
	}
}

// ForgotPassword generates reset token and sends email
func ForgotPassword(db *sqlx.DB, cfg *config.Config) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req struct {
			Email string `json:"email"`
		}
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		token := uuid.NewString()
		expiry := time.Now().Add(1 * time.Hour)
		if _, err := db.Exec(
			"UPDATE users SET reset_token=?, token_expiry=? WHERE email=?",
			token, expiry, req.Email,
		); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		link := fmt.Sprintf("%s/reset-password?token=%s", cfg.ServerAddress, token)
		if err := utils.SendEmail(cfg, req.Email, "Password Reset", fmt.Sprintf("Click here to reset: %s", link)); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to send email")
		}
		return c.JSON(http.StatusOK, echo.Map{"message": "Reset link sent"})
	}
}

// ResetPassword validates token and updates password
func ResetPassword(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req struct {
			Token    string `json:"token"`
			Password string `json:"password"`
		}
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		var user models.User
		if err := db.Get(&user, "SELECT * FROM users WHERE reset_token=?", req.Token); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid token")
		}
		if time.Now().After(user.TokenExpiry) {
			return echo.NewHTTPError(http.StatusBadRequest, "Token expired")
		}
		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to hash password")
		}
		if _, err := db.Exec(
			"UPDATE users SET password=?, reset_token=NULL, token_expiry=NULL WHERE id=?",
			string(hashed), user.ID,
		); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, echo.Map{"message": "Password reset successful"})
	}
}
