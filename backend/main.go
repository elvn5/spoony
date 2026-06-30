package main

import (
	"log"
	"net/http"

	"tma-boilerplate/config"
	"tma-boilerplate/database"
	"tma-boilerplate/handlers"
	"tma-boilerplate/middleware"

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
	auth := api.Group("/auth")
	{
		auth.POST("/telegram-login", handlers.TelegramLogin)
		auth.POST("/logout", handlers.Logout)
		auth.GET("/me", middleware.AuthRequired(), handlers.GetMe)
		auth.PUT("/profile", middleware.AuthRequired(), handlers.UpdateProfile)
	}

	// Telegram webhook (no auth — called by Telegram servers)
	api.POST("/webhook/telegram", handlers.HandleWebhook)
	api.GET("/webhook/info", handlers.GetWebhookInfo)

	// Admin API
	admin := r.Group("/admin")
	{
		adminAPI := admin.Group("/api", middleware.AdminAuth())
		{
			adminAPI.GET("/stats", handlers.AdminGetStats)
			adminAPI.GET("/users", handlers.AdminListUsers)
			adminAPI.DELETE("/users/:id", handlers.AdminDeleteUser)
		}
	}

	if err := handlers.RegisterWebhook(); err != nil {
		log.Printf("Warning: Telegram webhook registration failed: %v", err)
	}

	addr := ":" + config.App.Port
	log.Printf("Server starting on %s (env: %s)", addr, config.App.Env)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
