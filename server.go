package main

import (
	"PBD_backend_go/routes"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func runEnv() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
}

func main() {

	// Create a new fiber instance with custom config
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		// Call the next handler
		err := c.Next()
		// Check if we got an error
		if err != nil {
			// We had an error, do something with it
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Return from middleware
		return nil
	})
	// Add your routes
	app.Get("/", func(c *fiber.Ctx) error {
		// Simulate an error
		return fiber.ErrBadRequest
	})

	runEnv()
	routes.SetupRoutes(app)

	app.Listen(`localhost:` + os.Getenv(`GO_PORT`))
}
