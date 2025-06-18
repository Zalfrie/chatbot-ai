package config

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"os"
)

type Config struct {
	DBUser, DBPass, DBName, DBHost             string
	ServerAddress, JWTSecret                   string
	EmailHost, EmailPort, EmailUser, EmailPass string
}

func LoadConfig() *Config {
	return &Config{
		DBUser:        os.Getenv("DB_USER"),
		DBPass:        os.Getenv("DB_PASS"),
		DBName:        os.Getenv("DB_NAME"),
		DBHost:        os.Getenv("DB_HOST"),
		ServerAddress: os.Getenv("SERVER_ADDR"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		EmailHost:     os.Getenv("EMAIL_HOST"),
		EmailPort:     os.Getenv("EMAIL_PORT"),
		EmailUser:     os.Getenv("EMAIL_USER"),
		EmailPass:     os.Getenv("EMAIL_PASS"),
	}
}

func InitDB(cfg *Config) *sqlx.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBName)
	db := sqlx.MustConnect("mysql", dsn)
	return db
}
