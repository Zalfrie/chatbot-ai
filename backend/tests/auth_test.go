package tests

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/zalfrie/chatbot-ai/backend/controllers"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRegisterLogin(t *testing.T) {
	db, _ := sqlx.Connect("mysql", "chatbot:pass@tcp(127.0.0.1:3306)/chatbot_ai")
	e := echo.New()

	// Register
	req := httptest.NewRequest(http.MethodPost, "/api/register", strings.NewReader(`{"name":"Test","email":"test@example.com","password":"secret"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if assert.NoError(t, controllers.Register(db)(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}

	// Login
	req = httptest.NewRequest(http.MethodPost, "/api/login", strings.NewReader(`{"email":"test@example.com","password":"secret"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	if assert.NoError(t, controllers.Login(db, "testsecret")(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}
