package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                string
	Env                 string
	LogLevel            string
	DatabaseURL         string
	TelegramBotToken    string
	TelegramBotUsername string
	WebhookURL          string
	NgrokAPIURL         string
	JWTSecret           string
	JWTExpirationHours  int
	CORSAllowedOrigins  string
	TelegramMiniAppURL  string
	AdminToken          string
}

var App *Config

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	jwtExpHours, _ := strconv.Atoi(getEnv("JWT_EXPIRATION_HOURS", "720"))

	App = &Config{
		Port:                getEnv("PORT", "8080"),
		Env:                 getEnv("ENV", "development"),
		LogLevel:            getEnv("LOG_LEVEL", "info"),
		DatabaseURL:         getEnv("DATABASE_URL", "postgres://tma_user:tma_password@localhost:5432/tma_boilerplate?sslmode=disable"),
		TelegramBotToken:    getEnv("TELEGRAM_BOT_TOKEN", ""),
		TelegramBotUsername: getEnv("TELEGRAM_BOT_USERNAME", ""),
		WebhookURL:          getEnv("WEBHOOK_URL", ""),
		NgrokAPIURL:         getEnv("NGROK_API_URL", ""),
		JWTSecret:           getEnv("JWT_SECRET", "change-this-in-production-min-32-chars"),
		JWTExpirationHours:  jwtExpHours,
		CORSAllowedOrigins:  getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:5173"),
		TelegramMiniAppURL:  getEnv("TELEGRAM_MINI_APP_URL", "http://localhost:5173"),
		AdminToken:          getEnv("ADMIN_TOKEN", ""),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
