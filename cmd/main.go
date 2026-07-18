package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/imrishabk/goauth/internal/database"
	"github.com/imrishabk/goauth/internal/generator"
	"github.com/imrishabk/goauth/internal/keys"
	"github.com/imrishabk/goauth/internal/router"
	"github.com/joho/godotenv"
)

func main() {
	if err := initializeApp(); err != nil {
		log.Fatalln("Failed to start service:", err)
		os.Exit(1)
	}
}

func initializeApp() error {
	// Load Environment Files
	godotenv.Load(".env")
	// Generate Keys
	err := keys.ConfigureKeys()
	if err != nil {
		return err
	}
	// Initialize Node number for the snowflake generator
	generator.InitializeNode()
	// Establish connection to our database.
	database.ConnectDB()
	// Initialize fiber app
	app := fiber.New()

	// Setup routes
	router.SetupRoutes(app)

	// Handle root routes
	app.Get("/", greetingResponse)
	app.Post("/", greetingResponse)

	// Serve the application
	app.Listen(":"+os.Getenv("PORT"), fiber.ListenConfig{EnablePrefork: false})

	return nil
}

func greetingResponse(c fiber.Ctx) error {
	return c.SendString("goauth is running\nMade with ❤️ by imrishabk.")
}
