package main

import (
	"log"
	"net/http"

	"spoony/config"
	"spoony/database"
	"spoony/features/admin"
	"spoony/features/auth"
	"spoony/features/news"
	"spoony/features/telegrambot"
	"spoony/features/trainer"
	"spoony/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	config.Load()
	database.Connect()
	database.RunMigrations()

	if config.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.Use(middleware.CORS())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := r.Group("/api")

	// Auth routes
	authGroup := api.Group("/auth")
	{
		authGroup.POST("/telegram-login", auth.TelegramLogin)
		authGroup.POST("/guest", auth.GuestLogin)
		authGroup.POST("/logout", auth.Logout)
		authGroup.GET("/me", middleware.AuthRequired(), auth.GetMe)
		authGroup.PUT("/profile", middleware.AuthRequired(), auth.UpdateProfile)
	}

	// Spoony learning content
	api.GET("/news", news.GetNews)

	learn := api.Group("", middleware.AuthRequired())
	{
		learn.GET("/levels", trainer.GetLevels)
		learn.GET("/levels/:id/cards", trainer.GetLevelCards)
		learn.GET("/levels/:id/theory", trainer.GetLevelTheory)
		learn.POST("/levels/:id/complete", trainer.CompleteLevel)
		learn.GET("/stats", trainer.GetUserStats)
	}

	// Telegram webhook (no auth — called by Telegram servers)
	api.POST("/webhook/telegram", telegrambot.HandleWebhook)
	api.GET("/webhook/info", telegrambot.GetWebhookInfo)
	api.GET("/telegram/bot-info", telegrambot.GetBotInfo)

	// Admin API
	adminGroup := r.Group("/admin")
	{
		adminAPI := adminGroup.Group("/api", admin.Auth())
		{
			adminAPI.GET("/stats", admin.AdminGetStats)
			adminAPI.GET("/users", admin.AdminListUsers)
			adminAPI.GET("/users/:id", admin.AdminGetUser)
			adminAPI.DELETE("/users/:id", admin.AdminDeleteUser)
			adminAPI.PUT("/users/:id/progress/:levelId", admin.AdminUpdateUserProgress)
			adminAPI.DELETE("/users/:id/progress/:levelId", admin.AdminResetUserProgress)

			adminAPI.GET("/news", admin.AdminListNews)
			adminAPI.POST("/news", admin.AdminCreateNews)
			adminAPI.PUT("/news/:id", admin.AdminUpdateNews)
			adminAPI.DELETE("/news/:id", admin.AdminDeleteNews)
		}
	}

	if err := telegrambot.RegisterWebhook(); err != nil {
		log.Printf("Warning: Telegram webhook registration failed: %v", err)
	}

	addr := ":" + config.App.Port
	log.Printf("Server starting on %s (env: %s)", addr, config.App.Env)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
