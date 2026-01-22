// Package router is used to configure routes
package router

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/nullrish/goauth/handler"
)

func SetupRoutes(app *fiber.App) {
	// Authentication Routes
	auth := app.Group("/api/auth", cors.New(cors.Config{
		AllowOrigins:     []string{"https://rishabkarki.com.np"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))
	auth.Post("/login", handler.Login)
	auth.Post("/register", handler.Register)
	auth.Post("/verify-auth", handler.VerifyAuth)
}
