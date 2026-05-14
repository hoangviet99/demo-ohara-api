package routes

import (
	"demodiqit_api/config"
	"demodiqit_api/controllers"
	"demodiqit_api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(rg *gin.RouterGroup, cfg *config.Config) {
	authController := controllers.NewAuthController(cfg)

	// --- Public routes: no authentication required ---
	auth := rg.Group("/auth")
	{
		auth.POST("/login", authController.Login)
	}

	// --- Protected routes: valid JWT required ---
	protected := rg.Group("/")
	protected.Use(middleware.JWTAuthMiddleware(cfg))
	{
		protected.GET("/me", authController.Me)
	}

	// --- Admin-only routes: valid JWT + "admin" role required ---
	admin := rg.Group("/admin")
	admin.Use(middleware.JWTAuthMiddleware(cfg))
	admin.Use(middleware.AdminOnlyMiddleware())
	{
		// Placeholder — admin routes will be added in future tasks
		admin.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "admin access confirmed"})
		})
	}
}

