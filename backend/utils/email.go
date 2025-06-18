package utils

import (
	"fmt"
	"github.com/zalfrie/chatbot-ai/backend/config"
	"net/smtp"
)

// SendEmail sends a simple HTML email
func SendEmail(cfg *config.Config, to, subject, body string) error {
	auth := smtp.PlainAuth("", cfg.EmailUser, cfg.EmailPass, cfg.EmailHost)
	msg := []byte(fmt.Sprintf("Subject: %s\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s", subject, body))
	addr := fmt.Sprintf("%s:%s", cfg.EmailHost, cfg.EmailPort)
	return smtp.SendMail(addr, auth, cfg.EmailUser, []string{to}, msg)
}
