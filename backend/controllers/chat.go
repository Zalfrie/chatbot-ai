package controllers

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/zalfrie/chatbot-ai/backend/models"
)

// upgrader configures the Upgrade from HTTP to WebSocket protocol
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// connected clients and broadcast channel
var (
	clients   = make(map[*websocket.Conn]int)
	broadcast = make(chan models.Message)
)

// MemoryList returns stored public messages
func MemoryList(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var msgs []models.Message
		err := db.Select(&msgs, "SELECT * FROM messages ORDER BY created_at DESC")
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, msgs)
	}
}

// DeleteMemory deletes a persisted message by ID (admin only)
func DeleteMemory(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if _, err := db.Exec("DELETE FROM messages WHERE id=?", id); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.NoContent(http.StatusNoContent)
	}
}

// WebSocketHandler upgrades HTTP to WebSocket and manages real-time chat
func WebSocketHandler(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}
		defer ws.Close()

		// register client
		userID := c.Get("user_id").(int)
		clients[ws] = userID

		// start broadcast listener
		go func() {
			for msg := range broadcast {
				for client := range clients {
					client.WriteJSON(msg)
				}
			}
		}()

		// read incoming messages
		for {
			var msg models.Message
			if err := ws.ReadJSON(&msg); err != nil {
				delete(clients, ws)
				break
			}
			msg.UserID = userID
			msg.CreatedAt = time.Now()

			// send to all clients
			broadcast <- msg

			// persist only public messages
			if !msg.Private {
				if _, err := db.NamedExec(
					`INSERT INTO messages (user_id, content, is_private, created_at)
                     VALUES (:user_id, :content, :is_private, :created_at)`,
					msg,
				); err != nil {
					// log error (omitted)
				}
			}
		}

		return nil
	}
}
