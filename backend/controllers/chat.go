package controllers

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/zalfrie/chatbot-ai/backend/models"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	clients   = make(map[*websocket.Conn]int)
	broadcast = make(chan models.Message)
)

// WebSocketHandler menggunakan Gorilla WebSocket
func WebSocketHandler(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}
		defer ws.Close()

		userID := c.Get("user_id").(int)
		clients[ws] = userID

		// Goroutine untuk broadcast ke semua client
		go func() {
			for msg := range broadcast {
				for client := range clients {
					client.WriteJSON(msg)
				}
			}
		}()

		// Loop baca pesan
		for {
			var msg models.Message
			if err := ws.ReadJSON(&msg); err != nil {
				delete(clients, ws)
				break
			}
			msg.UserID = userID
			msg.CreatedAt = time.Now()
			broadcast <- msg

			if !msg.Private {
				db.NamedExec(`INSERT INTO messages (user_id,content,is_private,created_at)
                    VALUES (:user_id,:content,:is_private,:created_at)`, msg)
			}
		}

		return nil
	}
}
